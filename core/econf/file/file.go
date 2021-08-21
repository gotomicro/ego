package file

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/econf/manager"
	"github.com/gotomicro/ego/core/elog"
)

// fileDataSource file provider.
type fileDataSource struct {
	path string
	//dir         string
	enableWatch bool
	changed     chan struct{}
	logger      *elog.Component
}

func init() {
	manager.Register(manager.DefaultScheme, &fileDataSource{})
}

// Parse ...
func (fp *fileDataSource) Parse(path string, watch bool) econf.ConfigType {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		elog.Panic("new datasource", elog.FieldErr(err))
	}
	//dir := xfile.CheckAndGetParentDir(absolutePath)
	fp.path = absolutePath
	//fp.dir = dir
	fp.enableWatch = watch
	fp.logger = elog.EgoLogger.With(elog.FieldComponent(econf.PackageName))

	if watch {
		fp.changed = make(chan struct{}, 1)
		go fp.watch()
	}

	return extParser(path)
}

func extParser(configAddr string) econf.ConfigType {
	ext := filepath.Ext(configAddr)
	switch ext {
	case ".json":
		return econf.ConfigTypeJSON
	case ".toml":
		return econf.ConfigTypeToml
	case ".yaml":
		return econf.ConfigTypeYaml
	default:
		elog.EgoLogger.Panic("data source: invalid configuration type")
	}
	return ""
}

// ReadConfig ...
func (fp *fileDataSource) ReadConfig() (content []byte, err error) {
	return ioutil.ReadFile(fp.path)
}

// Close ...
func (fp *fileDataSource) Close() error {
	close(fp.changed)
	return nil
}

// IsConfigChanged ...
func (fp *fileDataSource) IsConfigChanged() <-chan struct{} {
	return fp.changed
}

// Watch file and automate update.
func (fp *fileDataSource) watch() {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		fp.logger.Fatal("new file watcher", elog.FieldComponent("file datasource"), elog.FieldErr(err))
	}
	defer w.Close()

	configFile := filepath.Clean(fp.path)
	realConfigFile, _ := filepath.EvalSymlinks(fp.path)

	fp.logger.Info("read watch",
		elog.FieldComponent("file datasource"),
		elog.String("configFile", configFile),
		elog.String("realConfigFile", realConfigFile),
		//elog.String("dir", fp.fdir),
		elog.String("fppath", fp.path),
	)
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-w.Events:
				currentConfigFile, _ := filepath.EvalSymlinks(fp.path)

				fp.logger.Info("read watch event",
					elog.FieldComponent("file datasource"),
					elog.String("event", filepath.Clean(event.Name)),
					elog.String("path", filepath.Clean(fp.path)),
					elog.String("currentConfigFile", currentConfigFile),
					elog.String("realConfigFile", realConfigFile),
				)
				// we only care about the config file with the following cases:
				// 1 - if the config file was modified or created
				// 2 - if the real path to the config file changed (eg: k8s ConfigMap replacement)
				const writeOrCreateMask = fsnotify.Write | fsnotify.Create
				if (filepath.Clean(event.Name) == configFile &&
					event.Op&writeOrCreateMask != 0) ||
					(currentConfigFile != "" && currentConfigFile != realConfigFile) {
					realConfigFile = currentConfigFile
					fp.logger.Info("modified file", elog.FieldName(event.Name), elog.FieldAddr(realConfigFile))
					select {
					case fp.changed <- struct{}{}:
					default:
					}
				}
			case err := <-w.Errors:
				fp.logger.Error("read watch error", elog.FieldComponent("file datasource"), elog.FieldErr(err))
			}
		}
	}()
	err = w.Add(fp.path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

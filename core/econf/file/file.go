package file

import (
	"encoding/json"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/econf/manager"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/util/xfile"
	"github.com/gotomicro/ego/core/util/xgo"
)

// DataSourceFile defines file scheme
const DataSourceFile = "file"

// fileDataSource file provider.
type fileDataSource struct {
	path        string
	dir         string
	enableWatch bool
	changed     chan struct{}
	logger      *elog.Component
}

func Register() {
	manager.Register("file", &fileDataSource{})
}

func ExtParser(configAddr string) (econf.Unmarshaller, string) {
	ext := filepath.Ext(configAddr)
	switch ext {
	case ".json":
		return json.Unmarshal, "json"
	case ".toml":
		return toml.Unmarshal, "toml"
	default:
		// TODO 处理configAddr为ETCD的情况？
		elog.EgoLogger.Panic("data source: invalid configuration type")
	}
	return nil, ""
}

// Parse
func (fp *fileDataSource) Parse(path string, watch bool) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		elog.Panic("new datasource", elog.Any("err", err))
	}
	dir := xfile.CheckAndGetParentDir(absolutePath)
	fp.path = absolutePath
	fp.dir = dir
	fp.enableWatch = watch
	if watch {
		fp.changed = make(chan struct{}, 1)
		xgo.Go(fp.watch)
	}
	return
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
		elog.Fatal("new file watcher", elog.FieldComponent("file datasource"), elog.Any("err", err))
	}
	defer w.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-w.Events:
				elog.Debug("read watch event",
					elog.FieldComponent("file datasource"),
					elog.String("event", filepath.Clean(event.Name)),
					elog.String("path", filepath.Clean(fp.path)),
				)
				// we only care about the config file with the following cases:
				// 1 - if the config file was modified or created
				// 2 - if the real path to the config file changed
				const writeOrCreateMask = fsnotify.Write | fsnotify.Create
				if event.Op&writeOrCreateMask != 0 && filepath.Clean(event.Name) == filepath.Clean(fp.path) {
					log.Println("modified file: ", event.Name)
					select {
					case fp.changed <- struct{}{}:
					default:
					}
				}
			case err := <-w.Errors:
				// log.Println("error: ", err)
				elog.Error("read watch error", elog.FieldComponent("file datasource"), elog.Any("err", err))
			}
		}
	}()

	err = w.Add(fp.dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

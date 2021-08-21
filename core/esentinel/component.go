package esentinel

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	sentinelAPI "github.com/alibaba/sentinel-golang/api"
	sentinelConfig "github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/fsnotify/fsnotify"
	"github.com/gotomicro/ego/core/elog"
)

// PackageName 设置包名
const PackageName = "core.sentinel"

// Component 组件
type Component struct {
	//name   string
	//config *Config
	//logger *elog.Component
	//err    error
}

func newComponent(config *Config, logger *elog.Component) error {
	if config.FlowRulesFile != "" {
		_ = syncFlowRules(config.FlowRulesFile, logger)
		go watch(config.FlowRulesFile, logger)
	}
	configEntity := sentinelConfig.NewDefaultConfig()
	configEntity.Sentinel.App.Name = config.AppName
	configEntity.Sentinel.Log.Dir = config.LogPath
	return sentinelAPI.InitWithConfig(configEntity)
}

func syncFlowRules(filePath string, logger *elog.Component) (err error) {
	var rules []*flow.Rule
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Error("load sentinel flow rules", elog.FieldErr(err), elog.FieldKey(filePath))
		return err
	}

	if err := json.Unmarshal(content, &rules); err != nil {
		logger.Error("load sentinel flow rules", elog.FieldErr(err), elog.FieldKey(filePath))
		return err
	}
	if len(rules) > 0 {
		_, _ = flow.LoadRules(rules)
	}
	return nil
}

func watch(filePath string, logger *elog.Component) {
	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		elog.Panic("new datasource", elog.FieldErr(err))
	}
	//dir := xfile.CheckAndGetParentDir(absolutePath)
	w, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Panic("new sentinel file watcher", elog.FieldErr(err))
	}
	defer w.Close()

	configFile := filepath.Clean(absolutePath)
	realConfigFile, _ := filepath.EvalSymlinks(absolutePath)
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-w.Events:
				currentConfigFile, _ := filepath.EvalSymlinks(absolutePath)

				logger.Info("read sentinel watch event",
					elog.String("event", filepath.Clean(event.Name)),
					elog.String("path", filepath.Clean(absolutePath)),
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
					logger.Info("modified sentinel file", elog.FieldName(event.Name), elog.FieldAddr(realConfigFile))
					_ = syncFlowRules(realConfigFile, logger)
				}
			case err := <-w.Errors:
				logger.Error("read watch error", elog.FieldComponent("file datasource"), elog.FieldErr(err))
			}
		}
	}()
	err = w.Add(absolutePath)
	if err != nil {
		logger.Panic("dir err", elog.FieldErr(err))
	}
	<-done
}

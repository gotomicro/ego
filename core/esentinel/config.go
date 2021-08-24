package esentinel

import (
	"github.com/gotomicro/ego/core/eapp"
)

// Config 配置
type Config struct {
	AppName       string `json:"appName"`       // 应用名，默认从ego框架内部获取
	LogPath       string `json:"logPath"`       // 日志路径，默认./logs
	FlowRulesFile string `json:"flowRulesFile"` // 限流配置路径
}

// DefaultConfig returns default config for sentinel
func DefaultConfig() *Config {
	return &Config{
		AppName: eapp.Name(),
		LogPath: "./logs",
	}
}

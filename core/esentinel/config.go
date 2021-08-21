package esentinel

import (
	"github.com/gotomicro/ego/core/eapp"
)

// Config 配置
type Config struct {
	AppName       string `json:"appName"`
	LogPath       string `json:"logPath"`
	FlowRulesFile string `json:"flowRulesFile"`
}

// DefaultConfig returns default config for sentinel
func DefaultConfig() *Config {
	return &Config{
		AppName: eapp.Name(),
		LogPath: "./logs",
	}
}

package egin

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// Config HTTP config
type Config struct {
	Host                      string
	Port                      int
	Mode                      string
	DisableMetric             bool
	DisableTrace              bool
	SlowQueryThresholdInMilli int64
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Host:                      "0.0.0.0",
		Port:                      9090,
		Mode:                      gin.ReleaseMode,
		SlowQueryThresholdInMilli: 500, // 500ms
	}
}

// Address ...
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

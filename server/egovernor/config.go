package egovernor

import (
	"fmt"

	"github.com/gotomicro/ego/core/eflag"
)

// Config 配置
type Config struct {
	Host              string
	Port              int
	EnableLocalMainIP bool
	EnableConnTcp     bool
	Network           string
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Host:    eflag.String("host"),
		Network: "tcp4",
		Port:    9003,
	}
}

// Address 地址
func (config Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

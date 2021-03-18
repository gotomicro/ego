package egovernor

import (
	"fmt"

	"github.com/gotomicro/ego/core/util/xnet"
)

// Config 配置
type Config struct {
	Host    string
	Port    int
	Network string
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	host, port, err := xnet.GetLocalMainIP()
	if err != nil {
		host = "localhost"
	}

	return &Config{
		Host:    host,
		Network: "tcp4",
		Port:    port,
	}
}

// Address 地址
func (config Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

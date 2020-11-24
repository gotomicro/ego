package egovernor

import (
	"fmt"

	"github.com/gotomicro/ego/core/util/xnet"
)

// Config ...
type Config struct {
	Host    string
	Port    int
	Network string
	Enable  bool
}

// DefaultConfig represents default config
// User should construct config base on DefaultConfig
func DefaultConfig() *Config {
	host, port, err := xnet.GetLocalMainIP()
	if err != nil {
		host = "localhost"
	}

	return &Config{
		Enable:  true,
		Host:    host,
		Network: "tcp4",
		Port:    port,
	}
}

// Address ...
func (config Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

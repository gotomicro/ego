package egovernor

import (
	"fmt"

	"github.com/gotomicro/ego/core/util/xnet"
)

type Config struct {
	Host    string
	Port    int
	Network string
}

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

func (config Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

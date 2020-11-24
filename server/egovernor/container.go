package egovernor

import (
	"github.com/gotomicro/ego/core/conf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/flag"
)

type Container struct {
	config *Config
	name   string
	err    error
	logger *elog.Component
}

func DefaultContainer() *Container {
	return &Container{
		logger: elog.EgoLogger.With(elog.FieldMod("server.egovernor")),
	}
}

func Load(key string) *Container {
	c := DefaultContainer()
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		c.err = err
		return c
	}
	// 修改host信息
	if flag.String("host") != "" {
		config.Host = flag.String("host")
	}
	c.config = config
	c.name = key
	return c
}

func (c *Container) Build() *Component {
	return newComponent(c.name, c.config, c.logger)
}

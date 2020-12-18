package egovernor

import (
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/eflag"
	"github.com/gotomicro/ego/core/elog"
)

type Option func(c *Container)

type Container struct {
	config *Config
	name   string
	err    error
	logger *elog.Component
}

func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
		logger: elog.EgoLogger.With(elog.FieldComponent(PackageName)),
	}
}

func Load(key string) *Container {
	c := DefaultContainer()
	if err := econf.UnmarshalKey(key, &c.config); err != nil {
		c.err = err
		return c
	}
	// 修改host信息
	if eflag.String("host") != "" {
		c.config.Host = eflag.String("host")
	}
	c.logger = c.logger.With(elog.FieldComponentName(key))
	c.name = key
	return c
}

func WithHost(host string) Option {
	return func(c *Container) {
		c.config.Host = host
	}
}

func WithPort(port int) Option {
	return func(c *Container) {
		c.config.Port = port
	}
}

func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	return newComponent(c.name, c.config, c.logger)
}

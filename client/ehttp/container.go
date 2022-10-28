package ehttp

import (
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
)

// Option 选项
type Option func(c *Container)

// Container defines a component instance.
type Container struct {
	config *Config
	name   string
	logger *elog.Component
}

// DefaultContainer returns an default container.
func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
		logger: elog.EgoLogger.With(elog.FieldComponent(PackageName)),
	}
}

// Load 记载配置key
func Load(key string) *Container {
	c := DefaultContainer()
	c.logger = c.logger.With(elog.FieldComponentName(key))
	if err := econf.UnmarshalKey(key, &c.config); err != nil {
		c.logger.Panic("parse config error", elog.FieldErr(err), elog.FieldKey(key))
		return c
	}
	c.name = key
	return c
}

// Build constructs a specific component from container.
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}

	c.logger.With(elog.FieldAddr(c.config.Addr))
	return newComponent(c.name, c.config, c.logger)
}

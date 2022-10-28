package esentinel

import (
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
)

// Option overrides a Container's default configuration.
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

// Load parses container configuration from configuration provider, such as a toml file,
// then use the configuration to construct a component container.
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
func (c *Container) Build(options ...Option) {
	err := newComponent(c.config, c.logger)
	if err != nil {
		elog.Panic("sentinel build panic", elog.FieldErr(err))
	}
}

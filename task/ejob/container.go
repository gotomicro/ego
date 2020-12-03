package ejob

import "github.com/gotomicro/ego/core/elog"

type Option func(c *Container)

type Container struct {
	config *Config
	logger *elog.Component
}

func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
		logger: elog.EgoLogger.With(elog.FieldComponent(PackageName)),
	}
}

func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	return newComponent(c.config.Name, c.config, c.logger)
}

package ecron

import (
	"github.com/robfig/cron/v3"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
)

type Option func(c *Container)

type Container struct {
	config *Config
	name   string
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
		c.logger.Panic("parse config error", elog.FieldErr(err), elog.FieldKey(key))
		return c
	}
	c.logger = c.logger.With(elog.FieldComponentName(key))
	c.name = key
	return c
}

// Build ...
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}

	if c.config.EnableWithSeconds {
		c.config.parser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	}

	switch c.config.DelayExecType {
	case "skip":
		c.config.wrappers = append(c.config.wrappers, skipIfStillRunning(c.logger))
	case "queue":
		c.config.wrappers = append(c.config.wrappers, queueIfStillRunning(c.logger))
	case "concurrent":
	default:
		c.config.wrappers = append(c.config.wrappers, skipIfStillRunning(c.logger))
	}

	if c.config.EnableDistributedTask && c.config.locker == nil {
		c.logger.Panic("client locker nil", elog.FieldKey("use WithLocker method"))
	}

	return newComponent(c.name, c.config, c.logger)
}

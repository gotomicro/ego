package ecron

import (
	"strings"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
)

// Container 容器
type Container struct {
	config *Config
	name   string
	logger *elog.Component
}

// DefaultContainer 默认容器
func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
		logger: elog.EgoLogger.With(elog.FieldComponent(PackageName)),
	}
}

// Load 加载配置key
func Load(key string) *Container {
	c := DefaultContainer()
	if err := econf.UnmarshalKey(key, c.config); err != nil {
		c.logger.Panic("parse config error", elog.FieldErr(err), elog.FieldKey(key))
		return c
	}
	c.config.Spec = strings.TrimSpace(c.config.Spec)
	c.logger = c.logger.With(elog.FieldComponentName(key))
	c.name = key
	return c
}

// Build 构建组件
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}

	if c.config.EnableSeconds {
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

	if c.config.EnableDistributedTask && c.config.lock == nil {
		c.logger.Panic("lock can not be nil", elog.FieldKey("use WithLock option to set lock"))
	}

	_, err := c.config.parser.Parse(c.config.Spec)
	if err != nil {
		c.logger.Panic("invalid cron spec", zap.Error(err))
	}

	return newComponent(c.name, c.config, c.logger)
}

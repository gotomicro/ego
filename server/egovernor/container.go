package egovernor

import (
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/eflag"
	"github.com/gotomicro/ego/core/elog"
)

// Container 容器
type Container struct {
	config *Config
	name   string
	err    error
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
	c.logger = c.logger.With(elog.FieldComponentName(key))
	if err := econf.UnmarshalKey(key, &c.config); err != nil {
		c.err = err
		return c
	}
	// 修改host信息
	if eflag.String("host") != "" {
		c.config.Host = eflag.String("host")
	}
	c.name = key
	return c
}

// Build 构建组件
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	return newComponent(c.name, c.config, c.logger)
}

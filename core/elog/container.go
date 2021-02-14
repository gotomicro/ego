package elog

import (
	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/econf"
)

type Option func(c *Container)

type Container struct {
	Config *Config
	name   string
}

func DefaultContainer() *Container {
	return &Container{
		Config: DefaultConfig(),
	}
}

func Load(key string) *Container {
	c := DefaultContainer()
	if err := econf.UnmarshalKey(key, &c.Config); err != nil {
		panic(err)
	}
	c.name = key
	return c
}

func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}

	if eapp.IsDevelopmentMode() {
		c.Config.Debug = true
		c.Config.EnableAsync = false
	}

	if c.Config.encoderConfig == nil {
		c.Config.encoderConfig = DefaultZapConfig()
	}

	if c.Config.Debug {
		c.Config.encoderConfig = DefaultDebugConfig()
	}

	if eapp.EnableLoggerAddApp() {
		c.Config.fields = append(c.Config.fields, FieldApp(eapp.Name()))
	}

	logger := newLogger(c.name, c.Config)
	// 如果名字不为空，加载动态配置
	if c.name != "" {
		// c.name 为配置name
		logger.AutoLevel(c.name + ".level")
	}

	return logger
}

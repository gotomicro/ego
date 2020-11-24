package elog

import (
	"github.com/gotomicro/ego/core/app"
	"github.com/gotomicro/ego/core/conf"
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
	if err := conf.UnmarshalKey(key, &c.Config); err != nil {
		panic(err)
		return c
	}
	c.name = key
	return c
}

func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}

	if app.IsDevelopmentMode() {
		c.Config.Debug = true
		c.Config.Async = false
		c.Config.EncoderConfig = DefaultDebugConfig()
	}

	if c.Config.EncoderConfig == nil {
		c.Config.EncoderConfig = DefaultZapConfig()
	}
	if c.Config.Debug {
		c.Config.EncoderConfig.EncodeLevel = DebugEncodeLevel
	}
	logger := newLogger(c.name, c.Config)
	// 如果名字不为空，加载动态配置
	if c.name != "" {
		// c.name 为配置name
		logger.AutoLevel(c.name + ".level")
	}
	return logger
}

func WithFileName(name string) Option {
	return func(c *Container) {
		c.Config.Name = name
	}
}

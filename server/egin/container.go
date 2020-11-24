package egin

import (
	"github.com/gotomicro/ego/core/conf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/flag"
)

type Option func(c *Container)

type Container struct {
	config *Config
	name   string
	logger *elog.Component
}

func DefaultContainer() *Container {
	return &Container{
		logger: elog.EgoLogger.With(elog.FieldMod("server.egin")),
	}
}

func Load(key string) *Container {
	c := DefaultContainer()
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		c.logger.Panic("parse config error", elog.FieldErr(err), elog.FieldKey(key))
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

func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	server := newComponent(c.name, c.config, c.logger)
	server.Use(recoverMiddleware(c.logger, c.config.SlowQueryThresholdInMilli))

	if !c.config.DisableMetric {
		server.Use(metricServerInterceptor())
	}

	if !c.config.DisableTrace {
		server.Use(traceServerInterceptor())
	}
	return server
}

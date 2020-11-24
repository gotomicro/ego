package egrpc

import (
	"github.com/gotomicro/ego/core/conf"
	"github.com/gotomicro/ego/core/elog"
	"google.golang.org/grpc"
)

type Option func(c *Container)

type Container struct {
	config *Config
	name   string
	logger *elog.Component
}

func DefaultContainer() *Container {
	return &Container{
		logger: elog.EgoLogger.With(elog.FieldMod("client.egrpc")),
	}
}

func Load(key string) *Container {
	c := DefaultContainer()
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		c.logger.Panic("parse config error", elog.FieldErr(err), elog.FieldKey(key))
		return c
	}

	c.config = config
	c.name = key
	return c
}

// WithDialOption ...
func WithDialOption(opts ...grpc.DialOption) Option {
	return func(c *Container) {
		if c.config.dialOptions == nil {
			c.config.dialOptions = make([]grpc.DialOption, 0)
		}
		c.config.dialOptions = append(c.config.dialOptions, opts...)
	}
}

// Build ...
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}

	if c.config.Debug {
		c.config.dialOptions = append(c.config.dialOptions,
			grpc.WithChainUnaryInterceptor(debugUnaryClientInterceptor(c.config.Address)),
		)
	}

	if !c.config.DisableAidInterceptor {
		c.config.dialOptions = append(c.config.dialOptions,
			grpc.WithChainUnaryInterceptor(aidUnaryClientInterceptor()),
		)
	}

	if !c.config.DisableTimeoutInterceptor {
		c.config.dialOptions = append(c.config.dialOptions,
			grpc.WithChainUnaryInterceptor(timeoutUnaryClientInterceptor(c.logger, c.config.ReadTimeout, c.config.SlowThreshold)),
		)
	}

	if !c.config.DisableTraceInterceptor {
		c.config.dialOptions = append(c.config.dialOptions,
			grpc.WithChainUnaryInterceptor(traceUnaryClientInterceptor()),
		)
	}

	if !c.config.DisableAccessInterceptor {
		c.config.dialOptions = append(c.config.dialOptions,
			grpc.WithChainUnaryInterceptor(loggerUnaryClientInterceptor(c.logger, c.config.Name, c.config.AccessInterceptorLevel)),
		)
	}

	if !c.config.DisableMetricInterceptor {
		c.config.dialOptions = append(c.config.dialOptions,
			grpc.WithChainUnaryInterceptor(metricUnaryClientInterceptor(c.config.Name)),
		)
	}

	c.logger.With(elog.FieldAddr(c.config.Address))
	return newComponent(c.name, c.config, c.logger)
}

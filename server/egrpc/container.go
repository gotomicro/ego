package egrpc

import (
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/eflag"
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
	// 修改host信息
	if eflag.String("host") != "" {
		c.config.Host = eflag.String("host")
	}
	c.logger = c.logger.With(elog.FieldComponentName(key))
	c.name = key
	return c
}

// WithServerOption inject server option to grpc server
// User should not inject interceptor option, which is recommend by WithStreamInterceptor
// and WithUnaryInterceptor
func WithServerOption(options ...grpc.ServerOption) Option {
	return func(c *Container) {
		if c.config.serverOptions == nil {
			c.config.serverOptions = make([]grpc.ServerOption, 0)
		}
		c.config.serverOptions = append(c.config.serverOptions, options...)
	}
}

// WithStreamInterceptor inject stream interceptors to server option
func WithStreamInterceptor(interceptors ...grpc.StreamServerInterceptor) Option {
	return func(c *Container) {
		if c.config.streamInterceptors == nil {
			c.config.streamInterceptors = make([]grpc.StreamServerInterceptor, 0)
		}

		c.config.streamInterceptors = append(c.config.streamInterceptors, interceptors...)
	}
}

// WithUnaryInterceptor inject unary interceptors to server option
func WithUnaryInterceptor(interceptors ...grpc.UnaryServerInterceptor) Option {
	return func(c *Container) {
		if c.config.unaryInterceptors == nil {
			c.config.unaryInterceptors = make([]grpc.UnaryServerInterceptor, 0)
		}
		c.config.unaryInterceptors = append(c.config.unaryInterceptors, interceptors...)
	}
}

// Build ...
func (c *Container) Build(options ...Option) *Component {
	if c.config.EnableTraceInterceptor {
		options = append(options, WithUnaryInterceptor(traceUnaryServerInterceptor))
		options = append(options, WithStreamInterceptor(traceStreamServerInterceptor))
	}

	if c.config.EnableMetricInterceptor {
		options = append(options, WithUnaryInterceptor(prometheusUnaryServerInterceptor))
		options = append(options, WithStreamInterceptor(prometheusStreamServerInterceptor))
	}

	for _, option := range options {
		option(c)
	}

	var streamInterceptors = append(
		[]grpc.StreamServerInterceptor{defaultStreamServerInterceptor(c.logger, c.config)},
		c.config.streamInterceptors...,
	)

	var unaryInterceptors = append(
		[]grpc.UnaryServerInterceptor{defaultUnaryServerInterceptor(c.logger, c.config)},
		c.config.unaryInterceptors...,
	)

	c.config.serverOptions = append(c.config.serverOptions,
		grpc.ChainStreamInterceptor(streamInterceptors...),
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
	)

	return newComponent(c.name, c.config, c.logger)
}

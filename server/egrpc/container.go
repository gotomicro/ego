package egrpc

import (
	"github.com/gotomicro/ego/core/conf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/flag"
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
		logger: elog.EgoLogger.With(elog.FieldMod("server.egrpc")),
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
func WithStreamInterceptor(intes ...grpc.StreamServerInterceptor) Option {
	return func(c *Container) {
		if c.config.streamInterceptors == nil {
			c.config.streamInterceptors = make([]grpc.StreamServerInterceptor, 0)
		}

		c.config.streamInterceptors = append(c.config.streamInterceptors, intes...)
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
	for _, option := range options {
		option(c)
	}

	if !c.config.DisableTrace {
		c.config.unaryInterceptors = append(c.config.unaryInterceptors, traceUnaryServerInterceptor)
		c.config.streamInterceptors = append(c.config.streamInterceptors, traceStreamServerInterceptor)
	}

	if !c.config.DisableMetric {
		c.config.unaryInterceptors = append(c.config.unaryInterceptors, prometheusUnaryServerInterceptor)
		c.config.streamInterceptors = append(c.config.streamInterceptors, prometheusStreamServerInterceptor)
	}

	return newComponent(c.name, c.config, c.logger)
}

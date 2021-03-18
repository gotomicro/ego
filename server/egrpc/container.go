package egrpc

import (
	"google.golang.org/grpc"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/eflag"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/util/xnet"
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
	var (
		host string
		err  error
	)
	// 获取网卡ip
	if c.config.EnableLocalMainIP {
		host, _, err = xnet.GetLocalMainIP()
		if err != nil {
			host = ""
		}
		c.config.Host = host
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
	var streamInterceptors []grpc.StreamServerInterceptor
	var unaryInterceptors []grpc.UnaryServerInterceptor
	// trace 必须在最外层，否则无法取到trace信息，传递到其他中间件
	if c.config.EnableTraceInterceptor {
		unaryInterceptors = []grpc.UnaryServerInterceptor{traceUnaryServerInterceptor, defaultUnaryServerInterceptor(c.logger, c.config)}
		streamInterceptors = []grpc.StreamServerInterceptor{traceStreamServerInterceptor, defaultStreamServerInterceptor(c.logger, c.config)}
	} else {
		unaryInterceptors = []grpc.UnaryServerInterceptor{defaultUnaryServerInterceptor(c.logger, c.config)}
		streamInterceptors = []grpc.StreamServerInterceptor{defaultStreamServerInterceptor(c.logger, c.config)}
	}

	if c.config.EnableMetricInterceptor {
		options = append(options, WithUnaryInterceptor(prometheusUnaryServerInterceptor))
		options = append(options, WithStreamInterceptor(prometheusStreamServerInterceptor))
	}

	for _, option := range options {
		option(c)
	}

	streamInterceptors = append(
		streamInterceptors,
		c.config.streamInterceptors...,
	)

	unaryInterceptors = append(
		unaryInterceptors,
		c.config.unaryInterceptors...,
	)

	c.config.serverOptions = append(c.config.serverOptions,
		grpc.ChainStreamInterceptor(streamInterceptors...),
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
	)

	return newComponent(c.name, c.config, c.logger)
}

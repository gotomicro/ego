package egrpc

import (
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/util/xnet"
	"google.golang.org/grpc"
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
	c.logger = c.logger.With(elog.FieldComponentName(key))
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
	c.name = key
	return c
}

// Build ...
func (c *Container) Build(options ...Option) *Component {
	var streamInterceptors []grpc.StreamServerInterceptor
	var unaryInterceptors []grpc.UnaryServerInterceptor
	// trace 必须在最外层，否则无法取到trace信息，传递到其他中间件
	if c.config.EnableTraceInterceptor {
		unaryInterceptors = []grpc.UnaryServerInterceptor{traceUnaryServerInterceptor, c.defaultUnaryServerInterceptor()}
		streamInterceptors = []grpc.StreamServerInterceptor{traceStreamServerInterceptor, c.defaultStreamServerInterceptor()}
	} else {
		unaryInterceptors = []grpc.UnaryServerInterceptor{c.defaultUnaryServerInterceptor()}
		streamInterceptors = []grpc.StreamServerInterceptor{c.defaultStreamServerInterceptor()}
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

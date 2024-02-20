package egrpc

import (
	"google.golang.org/grpc"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
)

// Option overrides a Container's default configuration.
type Option func(c *Container)

// Container defines a component instance.
type Container struct {
	config *Config
	name   string
	logger *elog.Component
}

// DefaultContainer returns an default container.
func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
		logger: elog.EgoLogger.With(elog.FieldComponent(PackageName)),
	}
}

// Load parses container configuration from configuration provider, such as a toml file,
// then use the configuration to construct a component container.
func Load(key string) *Container {
	c := DefaultContainer()
	c.logger = c.logger.With(elog.FieldComponentName(key))
	if err := econf.UnmarshalKey(key, &c.config); err != nil {
		c.logger.Panic("parse config error", elog.FieldErr(err), elog.FieldKey(key))
		return c
	}
	c.logger = c.logger.With(elog.FieldAddr(c.config.Addr))
	c.name = key
	return c
}

// Build constructs a specific component from container.
func (c *Container) Build(options ...Option) *Component {
	var unaryInterceptors []grpc.UnaryClientInterceptor
	var streamInterceptors []grpc.StreamClientInterceptor
	// 最先执行trace
	if c.config.EnableTraceInterceptor {
		unaryInterceptors = append(unaryInterceptors, c.traceUnaryClientInterceptor())
		streamInterceptors = append(streamInterceptors, c.traceStreamClientInterceptor())
	}
	// 默认日志
	unaryInterceptors = append(unaryInterceptors, c.loggerUnaryClientInterceptor())
	if eapp.IsDevelopmentMode() {
		unaryInterceptors = append(unaryInterceptors, c.debugUnaryClientInterceptor())
	}
	if c.config.EnableAppNameInterceptor {
		unaryInterceptors = append(unaryInterceptors, c.defaultUnaryClientInterceptor())
		streamInterceptors = append(streamInterceptors, c.defaultStreamClientInterceptor())
	}
	if c.config.EnableTimeoutInterceptor {
		unaryInterceptors = append(unaryInterceptors, c.timeoutUnaryClientInterceptor())
	}
	if c.config.EnableMetricInterceptor {
		unaryInterceptors = append(unaryInterceptors, c.metricUnaryClientInterceptor())
	}
	for _, option := range options {
		option(c)
	}
	c.config.dialOptions = append(c.config.dialOptions,
		grpc.WithChainStreamInterceptor(streamInterceptors...),
		grpc.WithChainUnaryInterceptor(unaryInterceptors...),
	)
	return newComponent(c.name, c.config, c.logger)
}

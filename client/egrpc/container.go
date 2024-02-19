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
	// 最先执行trace
	if c.config.EnableTraceInterceptor {
		options = append(options,
			WithDialOption(grpc.WithChainUnaryInterceptor(c.traceUnaryClientInterceptor())),
			WithDialOption(grpc.WithChainStreamInterceptor(c.traceStreamClientInterceptor())),
		)
	}

	// 其次执行，自定义header头，这样才能赋值到ctx里
	// options = append(options, WithDialOption(grpc.WithChainUnaryInterceptor(customHeader(transport.CustomContextKeys()))))

	// 默认日志
	options = append(options, WithDialOption(grpc.WithChainUnaryInterceptor(c.loggerUnaryClientInterceptor())))

	if eapp.IsDevelopmentMode() {
		options = append(options, WithDialOption(grpc.WithChainUnaryInterceptor(c.debugUnaryClientInterceptor())))
	}

	if c.config.EnableAppNameInterceptor {
		options = append(options, WithDialOption(grpc.WithChainUnaryInterceptor(c.defaultUnaryClientInterceptor())))
		options = append(options, WithDialOption(grpc.WithChainStreamInterceptor(c.defaultStreamClientInterceptor())))
	}

	if c.config.EnableTimeoutInterceptor {
		options = append(options, WithDialOption(grpc.WithChainUnaryInterceptor(c.timeoutUnaryClientInterceptor())))
	}

	if c.config.EnableMetricInterceptor {
		options = append(options,
			WithDialOption(grpc.WithChainUnaryInterceptor(c.metricUnaryClientInterceptor())),
		)
	}
	for _, option := range options {
		option(c)
	}
	// 兼容代码直接配置 grpc.MaxCallRecvMsgSize
	// 并保持配置文件高优先级
	if c.config.MaxCallRecvMsgSize != DefaultMaxCallRecvMsgSize {
		WithDialOption(grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(c.config.MaxCallRecvMsgSize)))(c)
	}
	return newComponent(c.name, c.config, c.logger)
}

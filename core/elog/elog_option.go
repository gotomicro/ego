package elog

import (
	"go.uber.org/zap/zapcore"
)

// Option overrides a Container's default configuration.
type Option func(c *Container)

// WithFileName 设置文件名
func WithFileName(name string) Option {
	return func(c *Container) {
		c.config.Name = name
	}
}

// WithDefaultFileName 设置默认的文件名,只有在配置的文件名不存在或者为兜底默认值的时候才会生效
func WithDefaultFileName(name string) Option {
	return func(c *Container) {
		// 只有当配置的文件名为空或者为兜底默认值的时候才会生效
		// 因为Container的默认值为DefaultLoggerName 如果配置了文件名，那么就不会使用默认的文件名了
		if c.config.Name == "" || c.config.Name == DefaultLoggerName {
			c.config.Name = name
		}
	}
}

// WithDebug 设置在命令行显示
func WithDebug(debug bool) Option {
	return func(c *Container) {
		c.config.Debug = debug
	}
}

// WithLevel 设置级别
func WithLevel(level string) Option {
	return func(c *Container) {
		c.config.Level = level
	}
}

// WithEnableAsync 是否异步执行，默认异步
func WithEnableAsync(enableAsync bool) Option {
	return func(c *Container) {
		c.config.EnableAsync = enableAsync
	}
}

// WithEnableAddCaller 是否添加行号，默认不添加行号
func WithEnableAddCaller(enableAddCaller bool) Option {
	return func(c *Container) {
		c.config.EnableAddCaller = enableAddCaller
	}
}

// WithZapCore 添加ZapCore
func WithZapCore(core zapcore.Core) Option {
	return func(c *Container) {
		c.config.core = core
	}
}

// WithEncoderConfig 加入encode config
func WithEncoderConfig(encoderConfig *zapcore.EncoderConfig) Option {
	return func(c *Container) {
		c.config.encoderConfig = encoderConfig
	}
}

// WithCallSkip 支持自定义调用层级
func WithCallSkip(callerSkip int) Option {
	return func(c *Container) {
		c.config.CallerSkip = callerSkip
	}
}

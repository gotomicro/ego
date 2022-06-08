package elog

import (
	"go.uber.org/zap/zapcore"
)

// Option 可选项
type Option func(c *Container)

// WithFileName 设置文件名
func WithFileName(name string) Option {
	return func(c *Container) {
		c.config.Name = name
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

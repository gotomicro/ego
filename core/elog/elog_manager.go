package elog

import (
	"io"

	"go.uber.org/zap/zapcore"
)

var (
	registry map[string]WriterBuilder
)

// WriterBuilder 根据key初始化writer
type WriterBuilder interface {
	// Build(key string, commonConfig *Config, lv zap.AtomicLevel) (zapcore.Core, CloseFunc)
	Build(key string, commonConfig *Config) Writer
	Scheme() string
}

// Writer 日志interface
type Writer interface {
	zapcore.Core
	io.Closer
}

// Close 关闭
func (c CloseFunc) Close() error {
	return c()
}

// CloseFunc should be called when the caller exits to clean up buffers.
type CloseFunc func() error

// Register registers a dataSource creator function to the registry
func Register(builder WriterBuilder) {
	registry[builder.Scheme()] = builder
}

// Provider 根据配置地址，创建数据源
func Provider(scheme string) WriterBuilder {
	logger, ok := registry[scheme]
	if !ok {
		panic("unsupported writer, error writer is: " + scheme)
	}
	return logger
}

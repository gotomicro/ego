package elog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	registry map[string]Writer
)

// Writer 根据key初始化writer
type Writer interface {
	Load(key string, commonConfig *Config, lv zap.AtomicLevel) (zapcore.Core, CloseFunc)
}

// CloseFunc should be called when the caller exits to clean up buffers.
type CloseFunc func() error

// Register registers a dataSource creator function to the registry
func Register(adapter string, creator Writer) {
	registry[adapter] = creator
}

// Provider 根据配置地址，创建数据源
func Provider(adapter string) Writer {
	logger, flag := registry[adapter]
	if !flag {
		panic("unsupported writer, error writer is: " + adapter)
	}
	return logger
}

package ejob

import (
	"context"
)

// Config ...
type Config struct {
	Name      string
	startFunc func(ctx context.Context) error
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Name:      "",
		startFunc: nil,
	}
}

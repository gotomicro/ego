package ehttp

import (
	"github.com/gotomicro/ego/core/util/xtime"
	"time"
)

type Config struct {
	Debug                        bool
	RawDebug                     bool
	Address                      string
	ReadTimeout                  time.Duration
	SlowLogThreshold             time.Duration
	EnableAccessInterceptor      bool
	EnableAccessInterceptorReply bool
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Debug:            false,
		SlowLogThreshold: xtime.Duration("500ms"),
		ReadTimeout:      xtime.Duration("2s"),
	}
}

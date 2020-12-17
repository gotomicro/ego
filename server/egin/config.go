package egin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/util/xtime"
	"time"
)

// Config HTTP config
type Config struct {
	Host                    string        // IP地址，默认127.0.0.1
	Port                    int           // PORT端口，默认9001
	Mode                    string        // gin的模式，默认是release模式
	EnableMetricInterceptor bool          // 是否开启监控，默认开启
	EnableTraceInterceptor  bool          // 是否开启链路追踪，默认开启
	SlowLogThreshold        time.Duration // 服务慢日志，默认500ms
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Host:                    "0.0.0.0",
		Port:                    9090,
		Mode:                    gin.ReleaseMode,
		EnableTraceInterceptor:  true,
		EnableMetricInterceptor: true,
		SlowLogThreshold:        xtime.Duration("500ms"),
	}
}

// Address ...
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

package egrpc

import (
	"fmt"
	"time"

	"google.golang.org/grpc"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/util/xtime"
)

// Config ...
type Config struct {
	Host                       string        // IP地址，默认0.0.0.0
	Port                       int           // Port端口，默认9002
	Deployment                 string        // 部署区域
	Network                    string        // 网络类型，默认tcp4
	EnableMetricInterceptor    bool          // 是否开启监控，默认开启
	EnableTraceInterceptor     bool          // 是否开启链路追踪，默认开启
	SlowLogThreshold           time.Duration // 服务慢日志，默认500ms
	EnableAccessInterceptorReq bool          // 是否开启记录请求参数，默认不开启
	EnableAccessInterceptorRes bool          // 是否开启记录响应参数，默认不开启

	serverOptions      []grpc.ServerOption
	streamInterceptors []grpc.StreamServerInterceptor
	unaryInterceptors  []grpc.UnaryServerInterceptor
}

// DefaultConfig represents default config
// User should construct config base on DefaultConfig
func DefaultConfig() *Config {
	return &Config{
		Network:                    "tcp4",
		Host:                       "0.0.0.0",
		Port:                       9002,
		Deployment:                 constant.DefaultDeployment,
		EnableMetricInterceptor:    true,
		EnableTraceInterceptor:     true,
		SlowLogThreshold:           xtime.Duration("500ms"),
		EnableAccessInterceptorReq: false,
		EnableAccessInterceptorRes: false,
		serverOptions:              []grpc.ServerOption{},
		streamInterceptors:         []grpc.StreamServerInterceptor{},
		unaryInterceptors:          []grpc.UnaryServerInterceptor{},
	}
}

// Address ...
func (config Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

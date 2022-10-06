package egrpc

import (
	"fmt"
	"time"

	"github.com/gotomicro/ego/core/eflag"

	"google.golang.org/grpc"

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
	EnableOfficialGrpcLog      bool          // 是否开启官方grpc日志，默认关闭
	EnableSkipHealthLog        bool          // 是否屏蔽探活日志，默认开启
	SlowLogThreshold           time.Duration // 服务慢日志，默认500ms
	EnableAccessInterceptor    bool          // 是否开启，记录请求数据
	EnableAccessInterceptorReq bool          // 是否开启记录请求参数，默认不开启
	EnableAccessInterceptorRes bool          // 是否开启记录响应参数，默认不开启
	EnableLocalMainIP          bool          // 自动获取ip地址
	serverOptions              []grpc.ServerOption
	streamInterceptors         []grpc.StreamServerInterceptor
	unaryInterceptors          []grpc.UnaryServerInterceptor
}

// DefaultConfig represents default config
// User should construct config base on DefaultConfig
func DefaultConfig() *Config {
	return &Config{
		Network:                    "tcp4",
		Host:                       eflag.String("host"),
		Port:                       9002,
		Deployment:                 "",
		EnableMetricInterceptor:    true,
		EnableSkipHealthLog:        true,
		EnableTraceInterceptor:     true,
		SlowLogThreshold:           xtime.Duration("500ms"),
		EnableAccessInterceptor:    true,
		EnableAccessInterceptorReq: false,
		EnableAccessInterceptorRes: false,
		serverOptions:              []grpc.ServerOption{},
		streamInterceptors:         []grpc.StreamServerInterceptor{},
		unaryInterceptors:          []grpc.UnaryServerInterceptor{},
	}
}

// Address ...
func (config Config) Address() string {
	// 如果是unix，那么启动方式为unix domain socket，host填写file
	if config.Network == "unix" {
		return config.Host
	}
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

package egrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/alibaba/sentinel-golang/core/base"
	"google.golang.org/grpc"

	"github.com/gotomicro/ego/core/eflag"
	"github.com/gotomicro/ego/core/util/xtime"
)

// Config ...
type Config struct {
	Host                          string        // IP地址，默认0.0.0.0
	Port                          int           // Port端口，默认9002
	Deployment                    string        // 部署区域
	Network                       string        // 网络类型，默认tcp4
	EnableMetricInterceptor       bool          // 是否开启监控，默认开启
	EnableTraceInterceptor        bool          // 是否开启链路追踪，默认开启
	EnableOfficialGrpcLog         bool          // 是否开启官方grpc日志，默认关闭
	EnableSkipHealthLog           bool          // 是否屏蔽探活日志，默认开启
	SlowLogThreshold              time.Duration // 服务慢日志，默认500ms
	EnableAccessInterceptor       bool          // 是否开启，记录请求数据
	EnableSentinel                bool          // 是否开启限流，默认不开启
	EnableAccessInterceptorReq    bool          // 是否开启记录请求参数，默认不开启
	AccessInterceptorReqMaxLength int           // 默认4K
	EnableAccessInterceptorRes    bool          // 是否开启记录响应参数，默认不开启
	AccessInterceptorResMaxLength int           // 默认4K
	EnableLocalMainIP             bool          // 自动获取ip地址
	serverOptions                 []grpc.ServerOption
	streamInterceptors            []grpc.StreamServerInterceptor
	unaryInterceptors             []grpc.UnaryServerInterceptor
	unaryServerResourceExtract    func(context.Context, interface{}, *grpc.UnaryServerInfo) string // sentinel 的限流策略
	unaryServerBlockFallback      func(context.Context, interface{}, *grpc.UnaryServerInfo, *base.BlockError) (interface{}, error)
}

// DefaultConfig represents default config
// User should construct config base on DefaultConfig
func DefaultConfig() *Config {
	return &Config{
		Network:                       "tcp4",
		Host:                          eflag.String("host"),
		Port:                          9002,
		Deployment:                    "",
		EnableMetricInterceptor:       true,
		EnableSkipHealthLog:           true,
		EnableTraceInterceptor:        true,
		EnableSentinel:                false,
		SlowLogThreshold:              xtime.Duration("500ms"),
		EnableAccessInterceptor:       true,
		EnableAccessInterceptorReq:    false,
		AccessInterceptorReqMaxLength: 4096,
		AccessInterceptorResMaxLength: 4096,
		EnableAccessInterceptorRes:    false,
		serverOptions:                 []grpc.ServerOption{},
		streamInterceptors:            []grpc.StreamServerInterceptor{},
		unaryInterceptors:             []grpc.UnaryServerInterceptor{},
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

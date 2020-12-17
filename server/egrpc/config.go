package egrpc

import (
	"fmt"
	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/util/xtime"
	"google.golang.org/grpc"
	"time"
)

// Config ...
type Config struct {
	Host                     string        // IP地址，默认0.0.0.0
	Port                     int           // Port端口，默认9002
	Deployment               string        // 部署区域
	Network                  string        // 网络类型，默认tcp4
	DisableTraceInterceptor  bool          // 禁用监控，默认否
	DisableMetricInterceptor bool          // 禁用trace，默认否
	SlowLogThreshold         time.Duration // 服务慢日志，默认500ms
	serverOptions            []grpc.ServerOption
	streamInterceptors       []grpc.StreamServerInterceptor
	unaryInterceptors        []grpc.UnaryServerInterceptor
}

// DefaultConfig represents default config
// User should construct config base on DefaultConfig
func DefaultConfig() *Config {
	return &Config{
		Network:                  "tcp4",
		Host:                     "0.0.0.0",
		Port:                     9002,
		Deployment:               constant.DefaultDeployment,
		DisableMetricInterceptor: false,
		DisableTraceInterceptor:  false,
		SlowLogThreshold:         xtime.Duration("500ms"),
		serverOptions:            []grpc.ServerOption{},
		streamInterceptors:       []grpc.StreamServerInterceptor{},
		unaryInterceptors:        []grpc.UnaryServerInterceptor{},
	}
}

// Address ...
func (config Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

package egrpc

import (
	"fmt"
	"github.com/gotomicro/ego/core/constant"
	"google.golang.org/grpc"
)

// Config ...
type Config struct {
	Host                    string // IP地址，默认0.0.0.0
	Port                    int    // Port端口，默认9002
	Deployment              string // 部署区域
	Network                 string // 网络类型，默认tcp4
	DisableTrace            bool   // 禁用监控，默认否
	DisableMetric           bool   // 禁用trace，默认否
	SlowLogThresholdInMilli int64  // 服务慢日志，默认500ms
	serverOptions           []grpc.ServerOption
	streamInterceptors      []grpc.StreamServerInterceptor
	unaryInterceptors       []grpc.UnaryServerInterceptor
}

// DefaultConfig represents default config
// User should construct config base on DefaultConfig
func DefaultConfig() *Config {
	return &Config{
		Network:                 "tcp4",
		Host:                    "0.0.0.0",
		Port:                    9002,
		Deployment:              constant.DefaultDeployment,
		DisableMetric:           false,
		DisableTrace:            false,
		SlowLogThresholdInMilli: 500,
		serverOptions:           []grpc.ServerOption{},
		streamInterceptors:      []grpc.StreamServerInterceptor{},
		unaryInterceptors:       []grpc.UnaryServerInterceptor{},
	}
}

// Address ...
func (config Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

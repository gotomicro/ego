package egrpc

import (
	"github.com/gotomicro/ego/core/util/xtime"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/keepalive"
)

// Config ...
type Config struct {
	BalancerName     string
	Addr             string
	Block            bool
	DialTimeout      time.Duration
	ReadTimeout      time.Duration
	Direct           bool
	OnFail           string // panic | error
	SlowLogThreshold time.Duration
	KeepAlive        *keepalive.ClientParameters
	dialOptions      []grpc.DialOption

	Debug                        bool
	DisableTraceInterceptor      bool
	DisableAppNameInterceptor    bool
	DisableTimeoutInterceptor    bool
	DisableMetricInterceptor     bool
	EnableAccessInterceptor      bool
	EnableAccessInterceptorReply bool
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		dialOptions: []grpc.DialOption{
			grpc.WithInsecure(),
		},
		Debug:            false,
		BalancerName:     roundrobin.Name, // round robin by default
		DialTimeout:      time.Second * 3,
		ReadTimeout:      xtime.Duration("1s"),
		SlowLogThreshold: xtime.Duration("600ms"),
		OnFail:           "panic",
		Block:            true,
	}
}

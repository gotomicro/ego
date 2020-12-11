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
	Name             string // config's name
	BalancerName     string
	Address          string
	Block            bool
	DialTimeout      time.Duration
	ReadTimeout      time.Duration
	Direct           bool
	OnDialError      string // panic | error
	KeepAlive        *keepalive.ClientParameters
	dialOptions      []grpc.DialOption
	SlowLogThreshold time.Duration

	Debug                     bool
	DisableTraceInterceptor   bool
	DisableAidInterceptor     bool
	DisableTimeoutInterceptor bool
	DisableMetricInterceptor  bool
	DisableAccessInterceptor  bool
	AccessInterceptorLevel    string
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		dialOptions: []grpc.DialOption{
			grpc.WithInsecure(),
		},
		BalancerName:           roundrobin.Name, // round robin by default
		DialTimeout:            time.Second * 3,
		ReadTimeout:            xtime.Duration("1s"),
		SlowLogThreshold:       xtime.Duration("600ms"),
		OnDialError:            "panic",
		AccessInterceptorLevel: "info",
		Block:                  true,
	}
}

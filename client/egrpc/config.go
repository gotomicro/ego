package egrpc

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/keepalive"

	"github.com/gotomicro/ego/core/util/xtime"
)

// Config ...
type Config struct {
	Addr                       string        // 连接地址，直连为127.0.0.1:9001，服务发现为etcd:///appname
	BalancerName               string        // 负载均衡方式，默认round robin
	OnFail                     string        // 失败后的处理方式，panic | error
	DialTimeout                time.Duration // 连接超时，默认3s
	ReadTimeout                time.Duration // 读超时，默认1s
	SlowLogThreshold           time.Duration // 慢日志记录的阈值，默认600ms
	Debug                      bool          // 是否开启调试，默认不开启，开启后并加上export EGO_DEBUG=true，可以看到每次请求，配置名、地址、耗时、请求数据、响应数据
	EnableBlock                bool          // 是否开启阻塞，默认开启
	EnableWithInsecure         bool          // 是否开启非安全传输，默认开启
	EnableMetricInterceptor    bool          // 是否开启监控，默认开启
	EnableTraceInterceptor     bool          // 是否开启链路追踪，默认开启
	EnableAppNameInterceptor   bool          // 是否开启传递应用名，默认开启
	EnableTimeoutInterceptor   bool          // 是否开启超时传递，默认开启
	EnableAccessInterceptor    bool          // 是否开启记录请求数据，默认不开启
	EnableAccessInterceptorReq bool          // 是否开启记录请求参数，默认不开启
	EnableAccessInterceptorRes bool          // 是否开启记录响应参数，默认不开启

	keepAlive   *keepalive.ClientParameters
	dialOptions []grpc.DialOption
}

// DefaultConfig defines grpc client default configuration
// User should construct config base on DefaultConfig
func DefaultConfig() *Config {
	return &Config{
		BalancerName:               roundrobin.Name,
		OnFail:                     "panic",
		DialTimeout:                time.Second * 3,
		ReadTimeout:                xtime.Duration("1s"),
		SlowLogThreshold:           xtime.Duration("600ms"),
		Debug:                      false,
		EnableBlock:                true,
		EnableTraceInterceptor:     true,
		EnableWithInsecure:         true,
		EnableAppNameInterceptor:   true,
		EnableTimeoutInterceptor:   true,
		EnableMetricInterceptor:    true,
		EnableAccessInterceptor:    false,
		EnableAccessInterceptorReq: false,
		EnableAccessInterceptorRes: false,
	}
}

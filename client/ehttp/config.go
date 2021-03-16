package ehttp

import (
	"github.com/gotomicro/ego/core/util/xtime"
	"runtime"
	"time"
)

type Config struct {
	Addr                       string        // 连接地址
	Debug                      bool          // 是否开启调试，默认不开启，开启后并加上export EGO_DEBUG=true，可以看到每次请求，配置名、地址、耗时、请求数据、响应数据
	RawDebug                   bool          // 是否开启原生调试，默认不开启
	ReadTimeout                time.Duration // 读超时，默认2s
	SlowLogThreshold           time.Duration // 慢日志记录的阈值，默认500ms
	IdleConnTimeout            time.Duration // 设置空闲连接时间，默认90 * time.Second
	MaxIdleConns               int           // 设置最大空闲连接数
	MaxIdleConnsPerHost        int           // 设置长连接个数
	EnableKeepAlives           bool          // 是否开启长连接，默认打开
	EnableAccessInterceptor    bool          // 是否开启记录请求数据，默认不开启
	EnableAccessInterceptorRes bool          // 是否开启记录响应参数，默认不开启
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Debug:                      false,
		RawDebug:                   false,
		ReadTimeout:                xtime.Duration("2s"),
		SlowLogThreshold:           xtime.Duration("500ms"),
		MaxIdleConns:               100,
		MaxIdleConnsPerHost:        runtime.GOMAXPROCS(0) + 1,
		IdleConnTimeout:            90 * time.Second,
		EnableKeepAlives:           true,
		EnableAccessInterceptor:    false,
		EnableAccessInterceptorRes: false,
	}
}

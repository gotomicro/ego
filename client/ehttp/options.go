package ehttp

import (
	"net/http"
	"time"
)

// WithAddr 设置HTTP地址
func WithAddr(addr string) Option {
	return func(c *Container) {
		c.config.Addr = addr
	}
}

// WithDebug 设置Debug信息
func WithDebug(debug bool) Option {
	return func(c *Container) {
		c.config.Debug = debug
	}
}

// WithRawDebug 设置原始Debug信息
func WithRawDebug(rawDebug bool) Option {
	return func(c *Container) {
		c.config.RawDebug = rawDebug
	}
}

// WithReadTimeout 设置读超时
func WithReadTimeout(readTimeout time.Duration) Option {
	return func(c *Container) {
		c.config.ReadTimeout = readTimeout
	}
}

// WithSlowLogThreshold 设置慢日志阈值
func WithSlowLogThreshold(slowLogThreshold time.Duration) Option {
	return func(c *Container) {
		c.config.SlowLogThreshold = slowLogThreshold
	}
}

// WithIdleConnTimeout 设置空闲连接时间
func WithIdleConnTimeout(idleConnTimeout time.Duration) Option {
	return func(c *Container) {
		c.config.IdleConnTimeout = idleConnTimeout
	}
}

// WithMaxIdleConns 设置最大空闲连接数
func WithMaxIdleConns(maxIdleConns int) Option {
	return func(c *Container) {
		c.config.MaxIdleConns = maxIdleConns
	}
}

// WithMaxIdleConnsPerHost 设置长连接个数
func WithMaxIdleConnsPerHost(maxIdleConnsPerHost int) Option {
	return func(c *Container) {
		c.config.MaxIdleConnsPerHost = maxIdleConnsPerHost
	}
}

// WithEnableTraceInterceptor 设置开启Trace拦截器
func WithEnableTraceInterceptor(enableTraceInterceptor bool) Option {
	return func(c *Container) {
		c.config.EnableTraceInterceptor = enableTraceInterceptor
	}
}

// WithEnableMetricInterceptor 设置开启 Metrics 拦截器
func WithEnableMetricInterceptor(enableMetricsInterceptor bool) Option {
	return func(c *Container) {
		c.config.EnableMetricInterceptor = enableMetricsInterceptor
	}
}

// WithEnableKeepAlives 设置是否开启长连接，默认打开
func WithEnableKeepAlives(enableKeepAlives bool) Option {
	return func(c *Container) {
		c.config.EnableKeepAlives = enableKeepAlives
	}
}

// WithEnableAccessInterceptor 设置开启请求日志
func WithEnableAccessInterceptor(enableAccessInterceptor bool) Option {
	return func(c *Container) {
		c.config.EnableAccessInterceptor = enableAccessInterceptor
	}
}

// WithEnableAccessInterceptorReq 设置开启请求日志响应参数
func WithEnableAccessInterceptorReq(enableAccessInterceptorReq bool) Option {
	return func(c *Container) {
		c.config.EnableAccessInterceptorReq = enableAccessInterceptorReq
	}
}

// WithEnableAccessInterceptorRes 设置开启请求日志响应参数
func WithEnableAccessInterceptorRes(enableAccessInterceptorRes bool) Option {
	return func(c *Container) {
		c.config.EnableAccessInterceptorRes = enableAccessInterceptorRes
	}
}

// WithPathRelabel 设置路径重命名
func WithPathRelabel(match string, replacement string) Option {
	return func(c *Container) {
		c.config.PathRelabel = append(c.config.PathRelabel, Relabel{Match: match, Replacement: replacement})
	}
}

// WithJar 设置Cookie，设置后，请求第一次接口后获取Cookie，第二次请求会带上Cookie，适合一些登录场景
// 例子：cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
func WithJar(jar http.CookieJar) Option {
	return func(c *Container) {
		c.config.cookieJar = jar
	}
}

// WithHTTPClient 设置自定义client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Container) {
		c.config.httpClient = httpClient
	}
}

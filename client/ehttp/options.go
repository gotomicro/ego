package ehttp

import "time"

func WithAddr(addr string) Option {
	return func(c *Container) {
		c.config.Addr = addr
	}
}

func WithDebug(debug bool) Option {
	return func(c *Container) {
		c.config.Debug = debug
	}
}

func WithRawDebug(rawDebug bool) Option {
	return func(c *Container) {
		c.config.RawDebug = rawDebug
	}
}

func WithReadTimeout(readTimeout time.Duration) Option {
	return func(c *Container) {
		c.config.ReadTimeout = readTimeout
	}
}

func WithSlowLogThreshold(slowLogThreshold time.Duration) Option {
	return func(c *Container) {
		c.config.SlowLogThreshold = slowLogThreshold
	}
}

func WithEnableAccessInterceptor(enableAccessInterceptor bool) Option {
	return func(c *Container) {
		c.config.EnableAccessInterceptor = enableAccessInterceptor
	}
}

func WithEnableAccessInterceptorRes(enableAccessInterceptorRes bool) Option {
	return func(c *Container) {
		c.config.EnableAccessInterceptorRes = enableAccessInterceptorRes
	}
}

// WithMaxIdleConns 设置最大空闲连接数
func WithMaxIdleConns(maxIdleConns int) Option {
	return func(c *Container) {
		c.config.MaxIdleConns = maxIdleConns
	}
}

// WithMaxIdleConns 设置长连接个数
func WithMaxIdleConnsPerHost(maxIdleConnsPerHost int) Option {
	return func(c *Container) {
		c.config.MaxIdleConnsPerHost = maxIdleConnsPerHost
	}
}

// WithEnableKeepAlives 设置是否开启长连接，默认打开
func WithEnableKeepAlives(enableKeepAlives bool) Option {
	return func(c *Container) {
		c.config.EnableKeepAlives = enableKeepAlives
	}
}

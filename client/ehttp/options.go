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

func WithEnableAccessInterceptorReply(enableAccessInterceptorReply bool) Option {
	return func(c *Container) {
		c.config.EnableAccessInterceptorReply = enableAccessInterceptorReply
	}
}

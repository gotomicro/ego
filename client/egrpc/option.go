package egrpc

import (
	"time"

	"google.golang.org/grpc"
)

// WithAddr setting grpc server address
func WithAddr(addr string) Option {
	return func(c *Container) {
		c.config.Addr = addr
	}
}

// WithOnFail setting failing mode
func WithOnFail(onFail string) Option {
	return func(c *Container) {
		c.config.OnFail = onFail
	}
}

// WithBalancerName setting grpc load balancer name
func WithBalancerName(balancerName string) Option {
	return func(c *Container) {
		c.config.BalancerName = balancerName
	}
}

// WithAddr setting grpc dial timeout
func WithDialTimeout(t time.Duration) Option {
	return func(c *Container) {
		c.config.DialTimeout = t
	}
}

// WithReadTimeout setting grpc read timeout
func WithReadTimeout(t time.Duration) Option {
	return func(c *Container) {
		c.config.ReadTimeout = t
	}
}

// WithDebug setting if enable debug mode
func WithDebug(enableDebug bool) Option {
	return func(c *Container) {
		c.config.Debug = enableDebug
	}
}

// WithDialOption setting grpc dial options
func WithDialOption(opts ...grpc.DialOption) Option {
	return func(c *Container) {
		if c.config.dialOptions == nil {
			c.config.dialOptions = make([]grpc.DialOption, 0)
		}
		c.config.dialOptions = append(c.config.dialOptions, opts...)
	}
}

package egrpc

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
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

// WithDialTimeout setting grpc dial timeout
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
		// for version compatibility
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

// WithEnableAccessInterceptor 开启日志记录
func WithEnableAccessInterceptor(enableAccessInterceptor bool) Option {
	return func(c *Container) {
		c.config.EnableAccessInterceptor = enableAccessInterceptor
	}
}

// WithEnableAccessInterceptorReq 开启日志请求参数
func WithEnableAccessInterceptorReq(enableAccessInterceptorReq bool) Option {
	return func(c *Container) {
		c.config.EnableAccessInterceptorReq = enableAccessInterceptorReq
	}
}

// WithEnableAccessInterceptorRes 开启日志响应记录
func WithEnableAccessInterceptorRes(enableAccessInterceptorRes bool) Option {
	return func(c *Container) {
		c.config.EnableAccessInterceptorRes = enableAccessInterceptorRes
	}
}

// WithBufnetServerListener 写入bufnet listener
func WithBufnetServerListener(svc net.Listener) Option {
	return WithDialOption(grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
		return svc.(*bufconn.Listener).Dial()
	}))
}

// WithName name
func WithName(name string) Option {
	return func(c *Container) {
		c.name = name
	}
}

package egrpc

import (
	"context"
	"google.golang.org/grpc"
)

// StreamInterceptorChain returns stream interceptors chain.
func StreamInterceptorChain(interceptors ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	build := func(c grpc.StreamServerInterceptor, n grpc.StreamHandler, info *grpc.StreamServerInfo) grpc.StreamHandler {
		return func(srv interface{}, stream grpc.ServerStream) error {
			return c(srv, stream, info, n)
		}
	}
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = build(interceptors[i], chain, info)
		}
		return chain(srv, stream)
	}
}

// UnaryInterceptorChain returns interceptors chain.
func UnaryInterceptorChain(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	build := func(c grpc.UnaryServerInterceptor, n grpc.UnaryHandler, info *grpc.UnaryServerInfo) grpc.UnaryHandler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			return c(ctx, req, info, n)
		}
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = build(interceptors[i], chain, info)
		}
		return chain(ctx, req)
	}
}

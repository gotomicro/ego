package egrpc

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"

	"github.com/gotomicro/ego/core/ecode"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/etrace"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/gotomicro/ego/core/emetric"
	"google.golang.org/grpc"
)

func prometheusUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	startTime := time.Now()
	resp, err := handler(ctx, req)
	code := ecode.ExtractCodes(err)
	emetric.ServerHandleHistogram.Observe(time.Since(startTime).Seconds(), emetric.TypeGRPCUnary, info.FullMethod, extractAID(ctx))
	emetric.ServerHandleCounter.Inc(emetric.TypeGRPCUnary, info.FullMethod, extractAID(ctx), code.GetMessage())
	return resp, err
}

func prometheusStreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	startTime := time.Now()
	err := handler(srv, ss)
	code := ecode.ExtractCodes(err)
	emetric.ServerHandleHistogram.Observe(time.Since(startTime).Seconds(), emetric.TypeGRPCStream, info.FullMethod, extractAID(ss.Context()))
	emetric.ServerHandleCounter.Inc(emetric.TypeGRPCStream, info.FullMethod, extractAID(ss.Context()), code.GetMessage())
	return err
}

func traceUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	span, ctx := etrace.StartSpanFromContext(
		ctx,
		info.FullMethod,
		etrace.FromIncomingContext(ctx),
		etrace.TagComponent("gRPC"),
		etrace.TagSpanKind("server.unary"),
	)

	defer span.Finish()

	resp, err := handler(ctx, req)

	if err != nil {
		code := codes.Unknown
		if s, ok := status.FromError(err); ok {
			code = s.Code()
		}
		span.SetTag("code", code)
		ext.Error.Set(span, true)
		span.LogFields(etrace.String("event", "error"), etrace.String("message", err.Error()))
	}
	return resp, err
}

type contextedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context ...
func (css contextedServerStream) Context() context.Context {
	return css.ctx
}

func traceStreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	span, ctx := etrace.StartSpanFromContext(
		ss.Context(),
		info.FullMethod,
		etrace.FromIncomingContext(ss.Context()),
		etrace.TagComponent("gRPC"),
		etrace.TagSpanKind("server.stream"),
		etrace.CustomTag("isServerStream", info.IsServerStream),
	)
	defer span.Finish()

	return handler(srv, contextedServerStream{
		ServerStream: ss,
		ctx:          ctx,
	})
}

func extractAID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return strings.Join(md.Get("aid"), ",")
	}
	return "unknown"
}

func defaultStreamServerInterceptor(logger *elog.Component, slowQueryThresholdInMilli int64) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		var beg = time.Now()
		var fields = make([]elog.Field, 0, 8)
		var event = "normal"
		defer func() {
			if slowQueryThresholdInMilli > 0 {
				if int64(time.Since(beg))/1e6 > slowQueryThresholdInMilli {
					event = "slow"
				}
			}

			if rec := recover(); rec != nil {
				switch rec := rec.(type) {
				case error:
					err = rec
				default:
					err = fmt.Errorf("%v", rec)
				}
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				fields = append(fields, elog.FieldStack(stack))
				event = "recover"
			}

			fields = append(fields,
				elog.Any("grpc interceptor type", "unary"),
				elog.FieldMethod(info.FullMethod),
				elog.FieldCost(time.Since(beg)),
				elog.FieldEvent(event),
			)

			for key, val := range getPeer(stream.Context()) {
				fields = append(fields, elog.Any(key, val))
			}

			if err != nil {
				fields = append(fields, zap.String("err", err.Error()))
				logger.Error("access", fields...)
				return
			}
			logger.Info("access", fields...)
		}()
		return handler(srv, stream)
	}
}

func defaultUnaryServerInterceptor(logger *elog.Component, slowQueryThresholdInMilli int64) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var beg = time.Now()
		var fields = make([]elog.Field, 0, 8)
		var event = "normal"
		defer func() {
			if slowQueryThresholdInMilli > 0 {
				if int64(time.Since(beg))/1e6 > slowQueryThresholdInMilli {
					event = "slow"
				}
			}
			if rec := recover(); rec != nil {
				switch rec := rec.(type) {
				case error:
					err = rec
				default:
					err = fmt.Errorf("%v", rec)
				}

				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				fields = append(fields, elog.FieldStack(stack))
				event = "recover"
			}

			fields = append(fields,
				elog.Any("grpc interceptor type", "unary"),
				elog.FieldMethod(info.FullMethod),
				elog.FieldCost(time.Since(beg)),
				elog.FieldEvent(event),
			)

			for key, val := range getPeer(ctx) {
				fields = append(fields, elog.Any(key, val))
			}

			if err != nil {
				fields = append(fields, zap.String("err", err.Error()))
				logger.Error("access", fields...)
				return
			}
			logger.Info("access", fields...)
		}()
		return handler(ctx, req)
	}
}

func getClientIP(ctx context.Context) (string, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("[getClinetIP] invoke FromContext() failed")
	}
	if pr.Addr == net.Addr(nil) {
		return "", fmt.Errorf("[getClientIP] peer.Addr is nil")
	}
	addSlice := strings.Split(pr.Addr.String(), ":")
	return addSlice[0], nil
}

func getPeer(ctx context.Context) map[string]string {
	var peerMeta = make(map[string]string)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val, ok := md["aid"]; ok {
			peerMeta["aid"] = strings.Join(val, ";")
		}
		var clientIP string
		if val, ok := md["client-ip"]; ok {
			clientIP = strings.Join(val, ";")
		} else {
			ip, err := getClientIP(ctx)
			if err == nil {
				clientIP = ip
			}
		}
		peerMeta["clientIP"] = clientIP
		if val, ok := md["client-host"]; ok {
			peerMeta["host"] = strings.Join(val, ";")
		}
	}
	return peerMeta

}

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

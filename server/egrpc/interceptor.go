package egrpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gotomicro/ego/core/eerrors"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/emetric"
	"github.com/gotomicro/ego/core/etrace"
	"github.com/gotomicro/ego/core/transport"
	"github.com/gotomicro/ego/core/util/xstring"
	"github.com/gotomicro/ego/internal/ecode"
	"github.com/gotomicro/ego/internal/tools"
	"github.com/gotomicro/ego/internal/xcpu"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func prometheusUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	startTime := time.Now()
	resp, err := handler(ctx, req)
	statusInfo := ecode.Convert(err)
	emetric.ServerHandleHistogram.Observe(time.Since(startTime).Seconds(), emetric.TypeGRPCUnary, info.FullMethod, getPeerName(ctx))
	emetric.ServerHandleCounter.Inc(emetric.TypeGRPCUnary, info.FullMethod, getPeerName(ctx), statusInfo.Message(), http.StatusText(ecode.GrpcToHTTPStatusCode(statusInfo.Code())))
	return resp, err
}

func prometheusStreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	startTime := time.Now()
	err := handler(srv, ss)
	statusInfo := ecode.Convert(err)
	emetric.ServerHandleHistogram.Observe(time.Since(startTime).Seconds(), emetric.TypeGRPCStream, info.FullMethod, getPeerName(ss.Context()))
	emetric.ServerHandleCounter.Inc(emetric.TypeGRPCUnary, info.FullMethod, getPeerName(ss.Context()), statusInfo.Message(), http.StatusText(ecode.GrpcToHTTPStatusCode(statusInfo.Code())))
	return err
}

func traceUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	tracer := etrace.NewTracer(trace.SpanKindServer)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		// Deprecated 该方法会在v0.9.0移除
		etrace.CompatibleExtractGrpcTraceId(md)
		ctx, span := tracer.Start(ctx, info.FullMethod, transport.GrpcHeaderCarrier(md))
		span.SetAttributes(
			etrace.TagComponent("grpc"),
			etrace.TagSpanKind("server.unary"),
		)
		defer func() {
			if err != nil {
				span.RecordError(err)
				if e := eerrors.FromError(err); e != nil {
					span.SetAttributes(attribute.Key("rpc.status_code").Int64(int64(e.Code)))
				}
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "OK")
			}
			span.End()
		}()
		return handler(ctx, req)
	}
}

type contextedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context ...
func (css contextedServerStream) Context() context.Context {
	return css.ctx
}

func traceStreamServerInterceptor() grpc.StreamServerInterceptor {
	tracer := etrace.NewTracer(trace.SpanKindServer)
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			md = metadata.New(nil)
		}
		// Deprecated 该方法会在v0.9.0移除
		etrace.CompatibleExtractGrpcTraceId(md)
		ctx, span := tracer.Start(ss.Context(), info.FullMethod, transport.GrpcHeaderCarrier(md))
		span.SetAttributes(
			etrace.TagComponent("grpc"),
			etrace.TagSpanKind("server.stream"),
		)
		defer span.End()
		return handler(srv, contextedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		})
	}
}

func (c *Container) defaultStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		var beg = time.Now()
		var fields = make([]elog.Field, 0, 20)
		var event = "normal"
		defer func() {
			cost := time.Since(beg)
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
			spbStatus := ecode.Convert(err)
			httpStatusCode := ecode.GrpcToHTTPStatusCode(spbStatus.Code())

			fields = append(fields,
				elog.FieldType("stream"),
				elog.FieldEvent(event),
				elog.FieldCode(int32(spbStatus.Code())),
				elog.FieldUniformCode(int32(httpStatusCode)),
				elog.FieldDescription(spbStatus.Message()),
				elog.FieldCost(time.Since(beg)),
				elog.FieldPeerName(getPeerName(stream.Context())),
				elog.FieldPeerIP(getPeerIP(stream.Context())),
			)

			if c.config.SlowLogThreshold > time.Duration(0) && c.config.SlowLogThreshold < cost {
				c.logger.Warn("slow", fields...)
			}

			if err != nil {
				fields = append(fields, elog.FieldErr(err))
				// 只记录系统级别错误
				if httpStatusCode >= http.StatusInternalServerError {
					// 只记录系统级别错误
					c.logger.Error("access", fields...)
				} else {
					// 非核心报错只做warning
					c.logger.Warn("access", fields...)
				}
				return
			}
			c.logger.Info("access", fields...)
		}()
		return handler(srv, stream)
	}
}

func (c *Container) defaultUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
		var beg = time.Now()
		// 为了性能考虑，如果要加日志字段，需要改变slice大小
		loggerKeys := transport.CustomContextKeys()
		var fields = make([]elog.Field, 0, 20+transport.CustomContextKeysLength())
		var event = "normal"

		// 必须在defer外层，因为要赋值，替换ctx
		for _, key := range loggerKeys {
			if value := tools.GrpcHeaderValue(ctx, key); value != "" {
				ctx = transport.WithValue(ctx, key, value)
			}
		}

		// 此处必须使用defer来recover handler内部可能出现的panic
		defer func() {
			cost := time.Since(beg)
			if rec := recover(); rec != nil {
				switch recType := rec.(type) {
				case error:
					err = recType
				default:
					err = fmt.Errorf("%v", rec)
				}

				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				fields = append(fields, elog.FieldStack(stack))
				event = "recover"
			}

			isSlow := false
			if c.config.SlowLogThreshold > time.Duration(0) && c.config.SlowLogThreshold < cost {
				isSlow = true
			}

			// 如果没有开启日志组件、并且没有错误，没有慢日志，那么直接返回不记录日志
			if err == nil && !c.config.EnableAccessInterceptor && !isSlow {
				return
			}

			spbStatus := ecode.Convert(err)
			httpStatusCode := ecode.GrpcToHTTPStatusCode(spbStatus.Code())

			fields = append(fields,
				elog.FieldType("unary"),
				elog.FieldCode(int32(spbStatus.Code())),
				elog.FieldUniformCode(int32(httpStatusCode)),
				elog.FieldDescription(spbStatus.Message()),
				elog.FieldEvent(event),
				elog.FieldMethod(info.FullMethod),
				elog.FieldCost(time.Since(beg)),
				elog.FieldPeerName(getPeerName(ctx)),
				elog.FieldPeerIP(getPeerIP(ctx)),
			)

			for _, key := range loggerKeys {
				if value := tools.ContextValue(ctx, key); value != "" {
					fields = append(fields, elog.FieldCustomKeyValue(key, value))
				}
			}

			if c.config.EnableTraceInterceptor && etrace.IsGlobalTracerRegistered() {
				fields = append(fields, elog.FieldTid(etrace.ExtractTraceID(ctx)))
			}

			if c.config.EnableAccessInterceptorReq {
				var reqMap = map[string]interface{}{
					"payload": xstring.JSON(req),
				}
				if md, ok := metadata.FromIncomingContext(ctx); ok {
					reqMap["metadata"] = md
				}
				fields = append(fields, elog.Any("req", reqMap))
			}
			if c.config.EnableAccessInterceptorRes {
				fields = append(fields, elog.Any("res", map[string]interface{}{
					"payload": xstring.JSON(res),
				}))
			}

			if isSlow {
				c.logger.Warn("slow", fields...)
			}

			if err != nil {
				fields = append(fields, elog.FieldErr(err))
				// 只记录系统级别错误
				if httpStatusCode >= http.StatusInternalServerError {
					// 只记录系统级别错误
					c.logger.Error("access", fields...)
				} else {
					// 非核心报错只做warning
					c.logger.Warn("access", fields...)
				}
				return
			}

			if c.config.EnableAccessInterceptor {
				c.logger.Info("access", fields...)
			}
		}()

		if enableCPUUsage(ctx) {
			var stat = xcpu.Stat{}
			xcpu.ReadStat(&stat)
			if stat.Usage > 0 {
				// https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md
				header := metadata.Pairs("cpu-usage", strconv.Itoa(int(stat.Usage)))
				err = grpc.SetHeader(ctx, header)
				if err != nil {
					c.logger.Error("set header error", elog.FieldErr(err))
				}
			}
		}
		return handler(ctx, req)
	}
}

// enableCPUUsage 是否开启cpu利用率
func enableCPUUsage(ctx context.Context) bool {
	return tools.GrpcHeaderValue(ctx, "enable-cpu-usage") == "true"
}

// getPeerName 获取对端应用名称
func getPeerName(ctx context.Context) string {
	return tools.GrpcHeaderValue(ctx, "app")
}

// getPeerIP 获取对端ip
func getPeerIP(ctx context.Context) string {
	clientIP := tools.GrpcHeaderValue(ctx, "client-ip")
	if clientIP != "" {
		return clientIP
	}

	// 从grpc里取对端ip
	pr, ok2 := peer.FromContext(ctx)
	if !ok2 {
		return ""
	}
	if pr.Addr == net.Addr(nil) {
		return ""
	}
	addSlice := strings.Split(pr.Addr.String(), ":")
	if len(addSlice) > 1 {
		return addSlice[0]
	}
	return ""
}

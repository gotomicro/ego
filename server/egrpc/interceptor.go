package egrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	sentinelbase "github.com/alibaba/sentinel-golang/core/base"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpccode "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/gotomicro/ego/core/eerrors"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/emetric"
	"github.com/gotomicro/ego/core/esentinel"
	"github.com/gotomicro/ego/core/etrace"
	"github.com/gotomicro/ego/core/transport"
	"github.com/gotomicro/ego/core/util/xstring"
	"github.com/gotomicro/ego/internal/ecode"
	"github.com/gotomicro/ego/internal/egrpcinteceptor"
	"github.com/gotomicro/ego/internal/tools"
)

const (
	mdKeyPeerName = "app"
	mdKeyPeerIp   = "client-ip"
)

func getPeerNameAndIp(ctx context.Context) (name string, ip string) {
	headers := tools.GrpcHeaderValues(ctx, mdKeyPeerName, mdKeyPeerIp)
	name = headers[0]
	ip = headers[1]
	if ip == "" {
		ip = getPeerIpFromContext(ctx)
	}
	return name, ip
}

func traceUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	tracer := etrace.NewTracer(trace.SpanKindServer)
	attrs := []attribute.KeyValue{
		egrpcinteceptor.RPCSystemGRPC,
		egrpcinteceptor.GRPCKindUnary,
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		// Deprecated 该方法会在v0.9.0移除
		// etrace.CompatibleExtractGrpcTraceID(md)
		ctx, span := tracer.Start(ctx, info.FullMethod, transport.GrpcHeaderCarrier(md), trace.WithAttributes(attrs...))
		peerName, peerIp := getPeerNameAndIp(ctx)
		span.SetAttributes(
			semconv.RPCMethodKey.String(info.FullMethod),
			semconv.NetPeerNameKey.String(peerName),
			semconv.NetPeerIPKey.String(peerIp),
		)
		defer func() {
			if err != nil {
				span.RecordError(err)
				if e := eerrors.FromError(err); e != nil {
					span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(e.Code)))
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

	receivedMessageID int
	sentMessageID     int
}

func (css *contextedServerStream) RecvMsg(m interface{}) error {
	err := css.ServerStream.RecvMsg(m)

	if err == nil {
		css.receivedMessageID++
		egrpcinteceptor.MessageReceived.Event(css.Context(), css.receivedMessageID, m)
	}

	return err
}

func (css *contextedServerStream) SendMsg(m interface{}) error {
	err := css.ServerStream.SendMsg(m)

	css.sentMessageID++
	egrpcinteceptor.MessageSent.Event(css.Context(), css.sentMessageID, m)

	return err
}

// Context ...
func (css *contextedServerStream) Context() context.Context {
	return css.ctx
}

func traceStreamServerInterceptor() grpc.StreamServerInterceptor {
	tracer := etrace.NewTracer(trace.SpanKindServer)
	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("grpc"),
		egrpcinteceptor.GRPCKindStream,
	}
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			md = metadata.New(nil)
		}
		// Deprecated 该方法会在v0.9.0移除
		// etrace.CompatibleExtractGrpcTraceID(md)
		ctx, span := tracer.Start(ss.Context(), info.FullMethod, transport.GrpcHeaderCarrier(md), trace.WithAttributes(attrs...))
		peerName, peerIp := getPeerNameAndIp(ctx)
		span.SetAttributes(
			semconv.RPCMethodKey.String(info.FullMethod),
			semconv.NetPeerNameKey.String(peerName),
			semconv.NetPeerIPKey.String(peerIp),
			etrace.CustomTag("rpc.grpc.kind", "stream"),
		)
		defer span.End()
		err := handler(srv, &contextedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		})
		if err != nil {
			span.RecordError(err)
			if e := eerrors.FromError(err); e != nil {
				span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(e.Code)))
			}
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "OK")
		}
		return err
	}
}

func (c *Container) defaultStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		var beg = time.Now()
		var fields []elog.Field
		var event = "normal"
		var spbStatus *status.Status
		var isSlowLog = false

		err = handler(srv, stream)
		cost := time.Since(beg)
		if c.config.SlowLogThreshold > time.Duration(0) && c.config.SlowLogThreshold < cost {
			event = "slow"
			isSlowLog = true
		}

		peerName, peerIp := getPeerNameAndIp(stream.Context())
		if c.config.EnableAccessInterceptor {
			fields = make([]elog.Field, 0, 20+transport.CustomContextKeysLength())
		}
		defer func() {
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
				err = status.New(grpccode.Internal, "panic recover, origin err: "+err.Error()).Err()
				fields = append(fields, elog.FieldKey("unary"), elog.FieldCode(int32(grpccode.Internal)),
					elog.FieldUniformCode(int32(http.StatusInternalServerError)), elog.FieldMethod(info.FullMethod),
					elog.FieldCost(time.Since(beg)), elog.FieldPeerName(peerName), elog.FieldType("recover"),
					elog.FieldPeerIP(peerIp), elog.FieldErr(err), elog.FieldStack(stack))
				c.logger.Error("access", fields...)
			}
		}()

		if c.config.EnableAccessInterceptor || err != nil || isSlowLog {
			spbStatus = ecode.Convert(err)
			httpStatusCode := ecode.GrpcToHTTPStatusCode(spbStatus.Code())
			fields = append(fields,
				elog.FieldKey("stream"),
				elog.FieldEvent(event),
				elog.FieldCode(int32(spbStatus.Code())),
				elog.FieldUniformCode(int32(httpStatusCode)),
				elog.FieldDescription(spbStatus.Message()),
				elog.FieldCost(time.Since(beg)),
				elog.FieldPeerName(peerName),
				elog.FieldPeerIP(peerIp),
			)

			if err != nil {
				// err!=nil, rpc处理报错时，记录额外的错误信息
				fields = append(fields, elog.FieldErr(err))
				// 只记录系统级别错误
				if httpStatusCode >= http.StatusInternalServerError {
					// 只记录系统级别错误
					c.logger.Error("access", fields...)
				} else {
					// 非核心报错只做warning
					c.logger.Warn("access", fields...)
				}
			} else if isSlowLog {
				// isSlowLog==true, 为慢日志时，记录日志
				c.logger.Warn("access", fields...)
			} else {
				// EnableAccessInterceptor==true, 开启了access日志，记录日志
				c.logger.Info("access", fields...)
			}
		}

		c.prometheusStreamServerInterceptor(stream, info, spbStatus, cost)
		return
	}
}

func (c *Container) prometheusStreamServerInterceptor(ss grpc.ServerStream, info *grpc.StreamServerInfo, pbStatus *status.Status, cost time.Duration) {
	serviceName, _ := egrpcinteceptor.SplitMethodName(info.FullMethod)
	emetric.ServerStartedCounter.Inc(emetric.TypeGRPCStream, info.FullMethod, getPeerName(ss.Context()), serviceName)
	// HandleHistogram的单位是s，需要用s单位
	emetric.ServerHandleHistogram.Observe(cost.Seconds(), emetric.TypeGRPCStream, info.FullMethod, getPeerName(ss.Context()), serviceName)
	emetric.ServerHandleCounter.Inc(emetric.TypeGRPCStream, info.FullMethod, getPeerName(ss.Context()), pbStatus.Message(), strconv.Itoa(ecode.GrpcToHTTPStatusCode(pbStatus.Code())), serviceName)
}

type ctxStore struct {
	mu  sync.RWMutex
	kvs map[string]any
}

type ctxStoreStruct struct{}

// CtxStoreSet 从ctx中尝试获取ctxStore，并往其中插入kv
func CtxStoreSet(ctx context.Context, k string, v any) {
	skv, ok := ctx.Value(ctxStoreStruct{}).(*ctxStore)
	if ok {
		skv.mu.Lock()
		skv.kvs[k] = v
		skv.mu.Unlock()
	}
}

func (c *Container) defaultUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
		ctx = context.WithValue(ctx, ctxStoreStruct{}, &ctxStore{kvs: map[string]any{}})
		// 默认过滤掉该探活日志
		if c.config.EnableSkipHealthLog && info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}

		var beg = time.Now()
		var fields []elog.Field
		var event = "normal"
		var spbStatus *status.Status
		var isSlowLog = false

		// 必须在defer外层，因为要赋值，替换ctx
		// 只有在环境变量里的自定义header，才会写入到context value里
		loggerKeys := transport.CustomContextKeys()
		headerKeys := append(loggerKeys, mdKeyPeerName, mdKeyPeerIp)
		headers := tools.GrpcHeaderValues(ctx, headerKeys...)
		for i, key := range headers[:len(headers)-2] {
			ctx = transport.WithValue(ctx, key, headers[i])
		}
		peerName := headers[len(headers)-2]
		peerIp := headers[len(headers)-1]
		if peerIp == "" {
			peerIp = getPeerIpFromContext(ctx)
		}
		if c.config.EnableAccessInterceptor {
			fields = make([]elog.Field, 0, 20+transport.CustomContextKeysLength())
		}

		// 此处必须使用defer来recover handler内部可能出现的panic
		defer func() {
			if rec := recover(); rec != nil {
				switch recType := rec.(type) {
				case error:
					err = recType
				default:
					err = fmt.Errorf("%v", rec)
				}

				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				err = status.New(grpccode.Internal, "panic recover, origin err: "+err.Error()).Err()
				fields = append(fields, elog.FieldKey("unary"), elog.FieldCode(int32(grpccode.Internal)),
					elog.FieldUniformCode(int32(http.StatusInternalServerError)), elog.FieldMethod(info.FullMethod),
					elog.FieldCost(time.Since(beg)), elog.FieldPeerName(peerName), elog.FieldType("recover"),
					elog.FieldPeerIP(peerIp), elog.FieldErr(err), elog.FieldStack(stack))
				c.logger.Error("access", fields...)
			}
		}()

		res, err = handler(ctx, req)
		cost := time.Since(beg)
		if c.config.SlowLogThreshold > time.Duration(0) && c.config.SlowLogThreshold < cost {
			isSlowLog = true
			event = "slow"
		}

		if c.config.EnableAccessInterceptor || err != nil || isSlowLog {
			spbStatus = ecode.Convert(err)
			httpStatusCode := ecode.GrpcToHTTPStatusCode(spbStatus.Code())
			fields = append(fields,
				elog.FieldKey("unary"),
				elog.FieldCode(int32(spbStatus.Code())),
				elog.FieldUniformCode(int32(httpStatusCode)),
				elog.FieldDescription(spbStatus.Message()),
				elog.FieldEvent(event),
				elog.FieldMethod(info.FullMethod),
				elog.FieldCost(time.Since(beg)),
				elog.FieldPeerName(peerName),
				elog.FieldPeerIP(peerIp),
			)

			skv, skvOk := ctx.Value(ctxStoreStruct{}).(*ctxStore)
			for _, key := range loggerKeys {
				if skvOk {
					skv.mu.RLock()
					v, ok := skv.kvs[key]
					skv.mu.RUnlock()
					if ok {
						fields = append(fields, elog.Any(strings.ToLower(key), v))
					}
				}
				if value := tools.ContextValue(ctx, key); value != "" {
					fields = append(fields, elog.FieldCustomKeyValue(key, value))
				}
			}

			if etrace.IsGlobalTracerRegistered() {
				fields = append(fields, elog.FieldTid(etrace.ExtractTraceID(ctx)))
			}

			if c.config.EnableAccessInterceptorReq {
				reqStr := xstring.JSON(req)
				if len(reqStr) > c.config.AccessInterceptorReqMaxLength {
					var reqMap = map[string]any{
						"payload": reqStr[:c.config.AccessInterceptorReqMaxLength] + "...",
					}
					if md, ok := metadata.FromIncomingContext(ctx); ok {
						reqMap["metadata"] = md
					}
					fields = append(fields, elog.Any("req", reqMap))
				} else {
					var reqMap = map[string]any{
						"payload": json.RawMessage(reqStr),
					}
					if md, ok := metadata.FromIncomingContext(ctx); ok {
						reqMap["metadata"] = md
					}
					fields = append(fields, elog.Any("req", reqMap))
				}
			}

			if c.config.EnableAccessInterceptorRes {
				resStr := xstring.JSON(res)
				if len(resStr) > c.config.AccessInterceptorResMaxLength {
					fields = append(fields, elog.Any("res", map[string]any{
						"payload": resStr[:c.config.AccessInterceptorResMaxLength] + "...",
					}))
				} else {
					fields = append(fields, elog.Any("res", map[string]any{
						"payload": json.RawMessage(resStr),
					}))
				}
			}

			if err != nil {
				// err!=nil, rpc处理报错时，记录额外的错误信息
				fields = append(fields, elog.FieldErr(err))
				// 只记录系统级别错误
				if httpStatusCode >= http.StatusInternalServerError {
					// 只记录系统级别错误
					c.logger.Error("access", fields...)
				} else {
					// 非核心报错只做warning
					c.logger.Warn("access", fields...)
				}
			} else if isSlowLog {
				// isSlowLog==true, 为慢日志时，记录日志
				c.logger.Warn("access", fields...)
			} else {
				// EnableAccessInterceptor==true, 开启了access日志，记录日志
				c.logger.Info("access", fields...)
			}
		}

		c.prometheusUnaryServerInterceptor(ctx, info, spbStatus, cost)
		return res, err
	}
}

func (c *Container) prometheusUnaryServerInterceptor(ctx context.Context, info *grpc.UnaryServerInfo, pbStatus *status.Status, cost time.Duration) {
	if !c.config.EnableMetricInterceptor {
		return
	}
	serviceName, _ := egrpcinteceptor.SplitMethodName(info.FullMethod)
	emetric.ServerStartedCounter.Inc(emetric.TypeGRPCUnary, info.FullMethod, getPeerName(ctx), serviceName)
	// HandleHistogram的单位是s，需要用s单位
	emetric.ServerHandleHistogram.ObserveWithExemplar(cost.Seconds(), prometheus.Labels{
		"tid": etrace.ExtractTraceID(ctx),
	}, emetric.TypeGRPCUnary, info.FullMethod, getPeerName(ctx), serviceName)
	emetric.ServerHandleCounter.Inc(emetric.TypeGRPCUnary, info.FullMethod, getPeerName(ctx), pbStatus.Code().String(), strconv.Itoa(ecode.GrpcToHTTPStatusCode(pbStatus.Code())), serviceName)
}

// getPeerName 获取对端应用名称
func getPeerName(ctx context.Context) string {
	return tools.GrpcHeaderValue(ctx, mdKeyPeerName)
}

// getPeerIP 获取对端ip
func getPeerIP(ctx context.Context) string {
	clientIP := tools.GrpcHeaderValue(ctx, mdKeyPeerIp)
	if clientIP != "" {
		return clientIP
	}

	return getPeerIpFromContext(ctx)
}

// getPeerIpFromContext 从grpc里取对端ip
func getPeerIpFromContext(ctx context.Context) string {
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

// NewUnaryServerInterceptor creates the unary server interceptor wrapped with Sentinel entry.
func (c *Container) sentinelInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// method as resource name by default
		resourceName := info.FullMethod
		if c.config.unaryServerResourceExtract != nil {
			resourceName = c.config.unaryServerResourceExtract(ctx, req, info)
		}

		if !esentinel.IsResExist(resourceName) {
			return handler(ctx, req)
		}

		// var entry *sentinelbase.SentinelEntry = nil
		entry, blockErr := sentinel.Entry(
			resourceName,
			sentinel.WithResourceType(sentinelbase.ResTypeRPC),
			sentinel.WithTrafficType(sentinelbase.Inbound),
		)
		if blockErr != nil {
			if c.config.unaryServerBlockFallback != nil {
				return c.config.unaryServerBlockFallback(ctx, req, info, blockErr)
			}

			return nil, eerrors.New(int(grpccode.ResourceExhausted), "blocked by sentinel", blockErr.Error())
		}
		defer entry.Exit()

		res, err := handler(ctx, req)
		if err != nil {
			sentinel.TraceError(entry, err)
		}
		return res, err
	}
}

package egrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/emetric"
	"github.com/gotomicro/ego/core/etrace"
	"github.com/gotomicro/ego/core/transport"
	"github.com/gotomicro/ego/core/util/xdebug"
	"github.com/gotomicro/ego/core/util/xstring"
	"github.com/gotomicro/ego/internal/ecode"
	"github.com/gotomicro/ego/internal/tools"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// metricUnaryClientInterceptor returns grpc unary request metrics collector interceptor
func (c *Container) metricUnaryClientInterceptor() func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		beg := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		statusInfo := ecode.Convert(err)
		emetric.ClientHandleCounter.Inc(emetric.TypeGRPCUnary, c.name, method, cc.Target(), statusInfo.Message())
		emetric.ClientHandleHistogram.Observe(time.Since(beg).Seconds(), emetric.TypeGRPCUnary, c.name, method, cc.Target())
		return err
	}
}

// debugUnaryClientInterceptor returns grpc unary request request and response details interceptor
func (c *Container) debugUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var p peer.Peer
		beg := time.Now()
		err := invoker(ctx, method, req, reply, cc, append(opts, grpc.Peer(&p))...)
		cost := time.Since(beg)
		if err != nil {
			log.Println("grpc.response", xdebug.MakeReqResErrorV2(6, c.name, c.config.Addr, cost, method+" | "+fmt.Sprintf("%v", req), err.Error()))
		} else {
			log.Println("grpc.response", xdebug.MakeReqResInfoV2(6, c.name, c.config.Addr, cost, method+" | "+fmt.Sprintf("%v", req), reply))
		}
		return err
	}
}

// traceUnaryClientInterceptor returns grpc unary opentracing interceptor
func traceUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}

		span, ctx := etrace.StartSpanFromContext(
			ctx,
			method,
			etrace.TagSpanKind("client"),
			etrace.TagComponent("grpc"),
		)
		defer span.Finish()

		err := invoker(etrace.MetadataInjector(ctx, md), method, req, reply, cc, opts...)
		if err != nil {
			span.SetTag("response_code", ecode.Convert(err).Code())
			ext.Error.Set(span, true)

			span.LogFields(etrace.String("event", "error"), etrace.String("message", err.Error()))
		}
		return err
	}
}

// defaultUnaryClientInterceptor returns interceptor inject app name
func (c *Container) defaultUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md
		ctx = metadata.AppendToOutgoingContext(ctx, "app", eapp.Name())
		if c.config.EnableCPUUsage {
			ctx = metadata.AppendToOutgoingContext(ctx, "enable-cpu-usage", "true")
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (c *Container) defaultStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		// https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md
		ctx = metadata.AppendToOutgoingContext(ctx, "app", eapp.Name())
		if c.config.EnableCPUUsage {
			ctx = metadata.AppendToOutgoingContext(ctx, "enable-cpu-usage", "true")
		}
		return streamer(ctx, desc, cc, method, opts...)
	}
}

// timeoutUnaryClientInterceptor settings timeout
func (c *Container) timeoutUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 若无自定义超时设置，默认设置超时
		_, ok := ctx.Deadline()
		if !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, c.config.ReadTimeout)
			defer cancel()
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// loggerUnaryClientInterceptor returns log interceptor for logging
func (c *Container) loggerUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, res interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		beg := time.Now()
		loggerKeys := transport.CustomContextKeys()
		var fields = make([]elog.Field, 0, 20+transport.CustomContextKeysLength())

		for _, key := range loggerKeys {
			if value := tools.ContextValue(ctx, key); value != "" {
				fields = append(fields, elog.FieldCustomKeyValue(key, value))
				// 替换context
				ctx = metadata.AppendToOutgoingContext(ctx, key, value)
			}
		}

		err := invoker(ctx, method, req, res, cc, opts...)
		cost := time.Since(beg)
		spbStatus := ecode.Convert(err)
		httpStatusCode := ecode.GrpcToHTTPStatusCode(spbStatus.Code())

		fields = append(fields,
			elog.FieldType("unary"),
			elog.FieldCode(int32(spbStatus.Code())),
			elog.FieldUniformCode(int32(httpStatusCode)),
			elog.FieldDescription(spbStatus.Message()),
			elog.FieldMethod(method),
			elog.FieldCost(cost),
			elog.FieldName(cc.Target()),
		)

		// 开启了链路，那么就记录链路id
		if c.config.EnableTraceInterceptor && opentracing.IsGlobalTracerRegistered() {
			fields = append(fields, elog.FieldTid(etrace.ExtractTraceID(ctx)))
		}

		if c.config.EnableAccessInterceptorReq {
			fields = append(fields, elog.Any("req", json.RawMessage(xstring.JSON(req))))
		}
		if c.config.EnableAccessInterceptorRes {
			fields = append(fields, elog.Any("res", json.RawMessage(xstring.JSON(res))))
		}

		if c.config.SlowLogThreshold > time.Duration(0) && cost > c.config.SlowLogThreshold {
			c.logger.Warn("slow", fields...)
		}

		if err != nil {
			fields = append(fields, elog.FieldEvent("error"), elog.FieldErr(err))
			// 只记录系统级别错误
			if httpStatusCode >= http.StatusInternalServerError {
				// 只记录系统级别错误
				c.logger.Error("access", fields...)
				return err
			}
			// 业务报错只做warning
			c.logger.Warn("access", fields...)
			return err
		}

		if c.config.EnableAccessInterceptor {
			fields = append(fields, elog.FieldEvent("normal"))
			c.logger.Info("access", fields...)
		}
		return nil
	}
}

// customHeader 自定义header头
func customHeader(egoLogExtraKeys []string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, res interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		for _, key := range egoLogExtraKeys {
			if value := tools.GrpcHeaderValue(ctx, key); value != "" {
				ctx = transport.WithValue(ctx, key, value)
			}
		}
		return invoker(ctx, method, req, res, cc, opts...)
	}
}

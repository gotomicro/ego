package egrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/ecode"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/emetric"
	"github.com/gotomicro/ego/core/etrace"
	"github.com/gotomicro/ego/core/util/xdebug"
	"github.com/gotomicro/ego/core/util/xstring"
)

// metricUnaryClientInterceptor returns grpc unary request metrics collector interceptor
func metricUnaryClientInterceptor(name string) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		beg := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)

		// 收敛err错误，将err过滤后，可以知道err是否为系统错误码
		spbStatus := ecode.ExtractCodes(err)
		// 只记录系统级别错误
		if spbStatus.Code < ecode.EcodeNum {
			// 只记录系统级别的详细错误码
			emetric.ClientHandleCounter.Inc(emetric.TypeGRPCUnary, name, method, cc.Target(), spbStatus.GetMessage())
			emetric.ClientHandleHistogram.Observe(time.Since(beg).Seconds(), emetric.TypeGRPCUnary, name, method, cc.Target())
		} else {
			emetric.ClientHandleCounter.Inc(emetric.TypeGRPCUnary, name, method, cc.Target(), "biz error")
			emetric.ClientHandleHistogram.Observe(time.Since(beg).Seconds(), emetric.TypeGRPCUnary, name, method, cc.Target())
		}
		return err
	}
}

// metricStreamClientInterceptor returns grpc stream request metrics collector interceptor
func metricStreamClientInterceptor(name string) func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		beg := time.Now()
		clientStream, err := streamer(ctx, desc, cc, method, opts...)

		// 暂时用默认的grpc的默认err收敛
		codes := ecode.ExtractCodes(err)
		emetric.ClientHandleCounter.Inc(emetric.TypeGRPCStream, name, method, cc.Target(), codes.GetMessage())
		emetric.ClientHandleHistogram.Observe(time.Since(beg).Seconds(), emetric.TypeGRPCStream, name, method, cc.Target())
		return clientStream, err
	}
}

// debugUnaryClientInterceptor returns grpc unary request request and response details interceptor
func debugUnaryClientInterceptor(logger *elog.Component, compName, addr string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var p peer.Peer
		prefix := fmt.Sprintf("[%s]", addr)
		if remote, ok := peer.FromContext(ctx); ok && remote.Addr != nil {
			prefix = prefix + "(" + remote.Addr.String() + ")"
		}

		beg := time.Now()
		err := invoker(ctx, method, req, reply, cc, append(opts, grpc.Peer(&p))...)
		cost := time.Since(beg)
		if eapp.IsDevelopmentMode() {
			if err != nil {
				log.Println("grpc.response", xdebug.MakeReqResError(compName, addr, cost, method+" | "+fmt.Sprintf("%v", req), err.Error()))
			} else {
				log.Println("grpc.response", xdebug.MakeReqResInfo(compName, addr, cost, method+" | "+fmt.Sprintf("%v", req), reply))
			}
		} else {
			// todo log
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
			code := codes.Unknown
			if s, ok := status.FromError(err); ok {
				code = s.Code()
			}
			span.SetTag("response_code", code)
			ext.Error.Set(span, true)

			span.LogFields(etrace.String("event", "error"), etrace.String("message", err.Error()))
		}
		return err
	}
}

// appNameUnaryClientInterceptor returns interceptor inject app name
func appNameUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		clientAppName := metadata.Pairs("app", eapp.Name())
		if ok {
			md = metadata.Join(md, clientAppName)
		} else {
			md = clientAppName
		}
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// timeoutUnaryClientInterceptor settings timeout
func timeoutUnaryClientInterceptor(_logger *elog.Component, timeout time.Duration, slowThreshold time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 若无自定义超时设置，默认设置超时
		_, ok := ctx.Deadline()
		if !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// loggerUnaryClientInterceptor returns log interceptor for logging
func loggerUnaryClientInterceptor(_logger *elog.Component, config *Config) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, res interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		beg := time.Now()
		err := invoker(ctx, method, req, res, cc, opts...)
		cost := time.Since(beg)
		isErrLog := false
		isSlowLog := false
		spbStatus := ecode.ExtractCodes(err)
		var fields = make([]elog.Field, 0, 15)
		fields = append(fields,
			elog.FieldType("unary"),
			elog.FieldCode(spbStatus.Code),
			elog.FieldDescription(spbStatus.Message),
			elog.FieldMethod(method),
			elog.FieldCost(cost),
			elog.FieldName(cc.Target()),
		)

		if config.EnableAccessInterceptorReq {
			fields = append(fields, elog.Any("req", json.RawMessage(xstring.Json(req))))
		}
		if config.EnableAccessInterceptorRes {
			fields = append(fields, elog.Any("res", json.RawMessage(xstring.Json(res))))
		}

		if err != nil {
			// 只记录系统级别错误
			if spbStatus.Code < ecode.EcodeNum {
				fields = append(fields, elog.FieldEvent("error"))
				// 只记录系统级别错误
				_logger.Error("access", fields...)
			} else {
				// 业务报错只做warning
				_logger.Warn("access", fields...)
			}
			isErrLog = true
			return err
		}

		if config.SlowLogThreshold > time.Duration(0) && cost > config.SlowLogThreshold {
			fields = append(fields, elog.FieldEvent("slow"))
			isSlowLog = true
			_logger.Warn("access", fields...)
		}

		if config.EnableAccessInterceptor && !isErrLog && !isSlowLog {
			fields = append(fields, elog.FieldEvent("normal"))
			_logger.Info("access", fields...)
		}
		return nil
	}
}

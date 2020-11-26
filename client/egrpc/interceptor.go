package egrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gotomicro/ego/core/app"
	"github.com/gotomicro/ego/core/ecode"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/etrace"
	"github.com/gotomicro/ego/core/metric"
	"github.com/gotomicro/ego/core/util/xcolor"
	"github.com/gotomicro/ego/core/util/xstring"
	"time"

	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var (
	errSlowCommand = errors.New("grpc unary slow command")
)

// metric统计
func metricUnaryClientInterceptor(name string) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		beg := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)

		// 收敛err错误，将err过滤后，可以知道err是否为系统错误码
		spbStatus := ecode.ExtractCodes(err)
		// 只记录系统级别错误
		if spbStatus.Code < ecode.EcodeNum {
			// 只记录系统级别的详细错误码
			metric.ClientHandleCounter.Inc(metric.TypeGRPCUnary, name, method, cc.Target(), spbStatus.GetMessage())
			metric.ClientHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeGRPCUnary, name, method, cc.Target())
		} else {
			metric.ClientHandleCounter.Inc(metric.TypeGRPCUnary, name, method, cc.Target(), "biz error")
			metric.ClientHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeGRPCUnary, name, method, cc.Target())
		}
		return err
	}
}

func metricStreamClientInterceptor(name string) func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		beg := time.Now()
		clientStream, err := streamer(ctx, desc, cc, method, opts...)

		// 暂时用默认的grpc的默认err收敛
		codes := ecode.ExtractCodes(err)
		metric.ClientHandleCounter.Inc(metric.TypeGRPCStream, name, method, cc.Target(), codes.GetMessage())
		metric.ClientHandleHistogram.Observe(time.Since(beg).Seconds(), metric.TypeGRPCStream, name, method, cc.Target())
		return clientStream, err
	}
}

func debugUnaryClientInterceptor(addr string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var p peer.Peer
		prefix := fmt.Sprintf("[%s]", addr)
		if remote, ok := peer.FromContext(ctx); ok && remote.Addr != nil {
			prefix = prefix + "(" + remote.Addr.String() + ")"
		}

		fmt.Printf("%-50s[%s] => %s\n", xcolor.Green(prefix), time.Now().Format("04:05.000"), xcolor.Green("Send: "+method+" | "+xstring.Json(req)))
		err := invoker(ctx, method, req, reply, cc, append(opts, grpc.Peer(&p))...)
		if err != nil {
			fmt.Printf("%-50s[%s] => %s\n", xcolor.Red(prefix), time.Now().Format("04:05.000"), xcolor.Red("Erro: "+err.Error()))
		} else {
			fmt.Printf("%-50s[%s] => %s\n", xcolor.Green(prefix), time.Now().Format("04:05.000"), xcolor.Green("Recv: "+xstring.Json(reply)))
		}

		return err
	}
}

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

func aidUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		clientAidMD := metadata.Pairs("aid", app.AppID())
		if ok {
			md = metadata.Join(md, clientAidMD)
		} else {
			md = clientAidMD
		}
		ctx = metadata.NewOutgoingContext(ctx, md)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// timeoutUnaryClientInterceptor gRPC客户端超时拦截器
func timeoutUnaryClientInterceptor(_logger *elog.Component, timeout time.Duration, slowThreshold time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		now := time.Now()
		// 若无自定义超时设置，默认设置超时
		_, ok := ctx.Deadline()
		if !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		err := invoker(ctx, method, req, reply, cc, opts...)
		du := time.Since(now)
		remoteIP := "unknown"
		if remote, ok := peer.FromContext(ctx); ok && remote.Addr != nil {
			remoteIP = remote.Addr.String()
		}

		if slowThreshold > time.Duration(0) && du > slowThreshold {
			_logger.Error("slow",
				elog.FieldErr(errSlowCommand),
				elog.FieldMethod(method),
				elog.FieldName(cc.Target()),
				elog.FieldCost(du),
				elog.FieldAddr(remoteIP),
			)
		}
		return err
	}
}

// loggerUnaryClientInterceptor gRPC客户端日志中间件
func loggerUnaryClientInterceptor(_logger *elog.Component, name string, accessInterceptorLevel string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		beg := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)

		spbStatus := ecode.ExtractCodes(err)
		if err != nil {
			// 只记录系统级别错误
			if spbStatus.Code < ecode.EcodeNum {
				// 只记录系统级别错误
				_logger.Error(
					"access",
					elog.FieldType("unary"),
					elog.FieldCode(spbStatus.Code),
					elog.FieldStringErr(spbStatus.Message),
					elog.FieldName(name),
					elog.FieldMethod(method),
					elog.FieldCost(time.Since(beg)),
					elog.Any("req", json.RawMessage(xstring.Json(req))),
					elog.Any("reply", json.RawMessage(xstring.Json(reply))),
				)
			} else {
				// 业务报错只做warning
				_logger.Warn(
					"access",
					elog.FieldType("unary"),
					elog.FieldCode(spbStatus.Code),
					elog.FieldStringErr(spbStatus.Message),
					elog.FieldName(name),
					elog.FieldMethod(method),
					elog.FieldCost(time.Since(beg)),
					elog.Any("req", json.RawMessage(xstring.Json(req))),
					elog.Any("reply", json.RawMessage(xstring.Json(reply))),
				)
			}
			return err
		} else {
			if accessInterceptorLevel == "info" {
				_logger.Info(
					"access",
					elog.FieldType("unary"),
					elog.FieldCode(spbStatus.Code),
					elog.FieldName(name),
					elog.FieldMethod(method),
					elog.FieldCost(time.Since(beg)),
					elog.Any("req", json.RawMessage(xstring.Json(req))),
					elog.Any("reply", json.RawMessage(xstring.Json(reply))),
				)
			}
		}

		return nil
	}
}

package egrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/proto"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/eerrors"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/emetric"
	"github.com/gotomicro/ego/core/etrace"
	"github.com/gotomicro/ego/core/transport"
	"github.com/gotomicro/ego/core/util/xdebug"
	"github.com/gotomicro/ego/core/util/xstring"
	"github.com/gotomicro/ego/internal/ecode"
	"github.com/gotomicro/ego/internal/egrpcinteceptor"
	"github.com/gotomicro/ego/internal/tools"
)

// metricUnaryClientInterceptor returns grpc unary request metrics collector interceptor
func (c *Container) metricUnaryClientInterceptor() func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		beg := time.Now()
		emetric.ClientStartedCounter.Inc(emetric.TypeGRPCUnary, c.name, method, cc.Target())
		err := invoker(ctx, method, req, reply, cc, opts...)
		statusInfo := ecode.Convert(err)

		emetric.ClientHandleHistogram.ObserveWithExemplar(time.Since(beg).Seconds(), prometheus.Labels{
			"tid": etrace.ExtractTraceID(ctx),
		}, emetric.TypeGRPCUnary, c.name, method, cc.Target())
		emetric.ClientHandleCounter.Inc(emetric.TypeGRPCUnary, c.name, method, cc.Target(), statusInfo.Code().String())
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
			log.Println("grpc.response", xdebug.MakeReqAndResError(fileWithLineNum(), c.name, c.config.Addr, cost, method+" | "+fmt.Sprintf("%v", req), err.Error()))
		} else {
			log.Println("grpc.response", xdebug.MakeReqAndResInfo(fileWithLineNum(), c.name, c.config.Addr, cost, method+" | "+fmt.Sprintf("%v", req), reply))
		}
		return err
	}
}

// traceUnaryClientInterceptor returns grpc unary opentracing interceptor
func (c *Container) traceUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	tracer := etrace.NewTracer(trace.SpanKindClient)
	attrs := []attribute.KeyValue{
		egrpcinteceptor.RPCSystemGRPC,
		egrpcinteceptor.GRPCKindUnary,
	}
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		ctx, span := tracer.Start(ctx, method, transport.GrpcHeaderCarrier(md), trace.WithAttributes(attrs...))
		span.SetAttributes(
			semconv.RPCMethodKey.String(method),
			semconv.NetPeerNameKey.String(c.config.Addr),
		)
		// 因为我们最先执行trace，所以这里，直接new出来metadata
		ctx = metadata.NewOutgoingContext(ctx, md)
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
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// defaultUnaryClientInterceptor returns interceptor inject app name
func (c *Container) defaultUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md
		ctx = metadata.AppendToOutgoingContext(ctx, "app", eapp.Name())
		// if c.config.EnableCPUUsage {
		//	ctx = metadata.AppendToOutgoingContext(ctx, "enable-cpu-usage", "true")
		// }
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

type streamEventType int

type streamEvent struct {
	Type streamEventType
	Err  error
}

const (
	receiveEndEvent streamEventType = iota
	errorEvent
)

var _ = proto.Marshal

// clientStream  wraps around the embedded grpc.ClientStream, and intercepts the RecvMsg and
// SendMsg method call.
type clientStream struct {
	grpc.ClientStream

	desc       *grpc.StreamDesc
	events     chan streamEvent
	eventsDone chan struct{}
	finished   chan error

	receivedMessageID int
	sentMessageID     int
}

func (w *clientStream) RecvMsg(m interface{}) error {
	err := w.ClientStream.RecvMsg(m)

	if err == nil && !w.desc.ServerStreams {
		w.sendStreamEvent(receiveEndEvent, nil)
	} else if err == io.EOF {
		w.sendStreamEvent(receiveEndEvent, nil)
	} else if err != nil {
		w.sendStreamEvent(errorEvent, err)
	} else {
		w.receivedMessageID++
		egrpcinteceptor.MessageReceived.Event(w.Context(), w.receivedMessageID, m)
	}

	return err
}

func (w *clientStream) SendMsg(m interface{}) error {
	err := w.ClientStream.SendMsg(m)

	w.sentMessageID++
	egrpcinteceptor.MessageSent.Event(w.Context(), w.sentMessageID, m)

	if err != nil {
		w.sendStreamEvent(errorEvent, err)
	}

	return err
}

func (w *clientStream) Header() (metadata.MD, error) {
	md, err := w.ClientStream.Header()

	if err != nil {
		w.sendStreamEvent(errorEvent, err)
	}

	return md, err
}

func (w *clientStream) CloseSend() error {
	err := w.ClientStream.CloseSend()

	if err != nil {
		w.sendStreamEvent(errorEvent, err)
	}

	return err
}

func wrapClientStream(ctx context.Context, s grpc.ClientStream, desc *grpc.StreamDesc) *clientStream {
	events := make(chan streamEvent)
	eventsDone := make(chan struct{})
	finished := make(chan error)

	go func() {
		defer close(eventsDone)

		for {
			select {
			case event := <-events:
				switch event.Type {
				case receiveEndEvent:
					finished <- nil
					return
				case errorEvent:
					finished <- event.Err
					return
				}
			case <-ctx.Done():
				finished <- ctx.Err()
				return
			}
		}
	}()

	return &clientStream{
		ClientStream: s,
		desc:         desc,
		events:       events,
		eventsDone:   eventsDone,
		finished:     finished,
	}
}

func (w *clientStream) sendStreamEvent(eventType streamEventType, err error) {
	select {
	case <-w.eventsDone:
	case w.events <- streamEvent{Type: eventType, Err: err}:
	}
}

func (c *Container) traceStreamClientInterceptor() grpc.StreamClientInterceptor {
	tracer := etrace.NewTracer(trace.SpanKindClient)
	attrs := []attribute.KeyValue{
		egrpcinteceptor.RPCSystemGRPC,
		egrpcinteceptor.GRPCKindStream,
	}
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		ctx, span := tracer.Start(ctx, method, transport.GrpcHeaderCarrier(md), trace.WithAttributes(attrs...))
		span.SetAttributes(
			semconv.RPCMethodKey.String(method),
			semconv.NetPeerNameKey.String(c.config.Addr),
		)
		// 因为我们最先执行trace，所以这里，直接new出来metadata
		ctx = metadata.NewOutgoingContext(ctx, md)

		s, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			span.RecordError(err)
			if e := eerrors.FromError(err); e != nil {
				span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(e.Code)))
			}
			span.SetStatus(codes.Error, err.Error())
			span.End()
			return s, err
		}
		stream := wrapClientStream(ctx, s, desc)

		go func() {
			err := <-stream.finished
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

		return stream, nil
	}
}

func (c *Container) defaultStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		// https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md
		ctx = metadata.AppendToOutgoingContext(ctx, "app", eapp.Name())
		// if c.config.EnableCPUUsage {
		//	ctx = metadata.AppendToOutgoingContext(ctx, "enable-cpu-usage", "true")
		// }
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
				// 替换context
				ctx = metadata.AppendToOutgoingContext(ctx, key, value)

				// grpc metadata 存在同一个 key set 多次，客户端可通过日志排查这种错误使用。
				if md, ok := metadata.FromOutgoingContext(ctx); ok {
					fields = append(fields, elog.FieldCustomKeyValue(key, strings.Join(md[strings.ToLower(key)], ";")))
				}
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
		if c.config.EnableTraceInterceptor && etrace.IsGlobalTracerRegistered() {
			fields = append(fields, elog.FieldTid(etrace.ExtractTraceID(ctx)))
		}

		if c.config.EnableAccessInterceptorReq {
			var reqMap = map[string]interface{}{
				"payload": xstring.JSON(req),
			}
			if md, ok := metadata.FromOutgoingContext(ctx); ok {
				reqMap["metadata"] = md
			}
			fields = append(fields, elog.Any("req", reqMap))
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

func fileWithLineNum() string {
	// the second caller usually from internal, so set i start from 2
	for i := 2; i < 20; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if (!(strings.Contains(file, "ego") && strings.HasSuffix(file, "client/egrpc/interceptor.go")) && !strings.HasSuffix(file, ".pb.go") && !strings.Contains(file, "google.golang.org")) || strings.HasSuffix(file, "_test.go") {
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}
	return ""
}

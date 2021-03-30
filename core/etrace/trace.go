package etrace

import (
	"context"

	"github.com/uber/jaeger-client-go"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/gotomicro/ego/core/elog"
)

var (
	// String ...
	String = log.String
)

// SetGlobalTracer ...
func SetGlobalTracer(tracer opentracing.Tracer) {
	elog.EgoLogger.Info("set global tracer", elog.FieldComponent("trace"))
	opentracing.SetGlobalTracer(tracer)
}

// StartSpanFromContext ...
func StartSpanFromContext(ctx context.Context, op string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, op, opts...)
}

// SpanFromContext ...
func SpanFromContext(ctx context.Context) opentracing.Span {
	return opentracing.SpanFromContext(ctx)
}

// ExtractTraceID HTTP使用request.Context，不要使用错了
func ExtractTraceID(ctx context.Context) string {
	if !opentracing.IsGlobalTracerRegistered() {
		return ""
	}

	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return ""
	}
	return span.(*jaeger.Span).Context().(jaeger.SpanContext).TraceID().String()
}

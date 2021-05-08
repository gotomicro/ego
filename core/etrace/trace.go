package etrace

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
)

var (
	// String ...
	String = log.String
)

// SetGlobalTracer ...
func SetGlobalTracer(tracer opentracing.Tracer) {
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

package etrace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var (
	// String ...
	String = attribute.String
)

type registeredTracer struct {
	isRegistered bool
}

var (
	globalTracer = registeredTracer{false}
)

// SetGlobalTracer ...
func SetGlobalTracer(tp trace.TracerProvider) {
	globalTracer = registeredTracer{true}
	otel.SetTracerProvider(tp)
}

// IsGlobalTracerRegistered returns a `bool` to indicate if a tracer has been globally registered
func IsGlobalTracerRegistered() bool {
	return globalTracer.isRegistered
}

// ExtractTraceID HTTP使用request.Context，不要使用错了
func ExtractTraceID(ctx context.Context) string {
	if !IsGlobalTracerRegistered() {
		return ""
	}
	span := trace.SpanContextFromContext(ctx)
	if span.HasTraceID() {
		return span.TraceID().String()
	}
	return ""
}

// Tracer is otel span tracer
type Tracer struct {
	tracer trace.Tracer
	kind   trace.SpanKind
	opt    *options
}

// NewTracer create tracer instance
func NewTracer(kind trace.SpanKind, opts ...Option) *Tracer {
	op := options{
		propagator: propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}),
	}
	for _, o := range opts {
		o(&op)
	}
	return &Tracer{tracer: otel.Tracer("ego"), kind: kind, opt: &op}
}

// Start start tracing span
func (t *Tracer) Start(ctx context.Context, operation string, carrier propagation.TextMapCarrier) (context.Context, trace.Span) {
	if (t.kind == trace.SpanKindServer || t.kind == trace.SpanKindConsumer) && carrier != nil {
		ctx = t.opt.propagator.Extract(ctx, carrier)
	}
	ctx, span := t.tracer.Start(ctx,
		operation,
		trace.WithSpanKind(t.kind),
	)
	if (t.kind == trace.SpanKindClient || t.kind == trace.SpanKindProducer) && carrier != nil {
		t.opt.propagator.Inject(ctx, carrier)
	}
	return ctx, span
}

type options struct {
	propagator propagation.TextMapPropagator
}

// Option is tracing option.
type Option func(*options)

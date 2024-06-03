package etrace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func TestExtractTraceID(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		globalTracer.isRegistered = false
		var ctx context.Context
		out := ExtractTraceID(ctx)
		assert.Equal(t, "", out)
	})

	t.Run("case 2", func(t *testing.T) {
		spanCtx := trace.NewSpanContext(trace.SpanContextConfig{TraceID: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}})
		ctx := trace.ContextWithSpanContext(context.Background(), spanCtx)
		globalTracer.isRegistered = true
		traceID := ExtractTraceID(ctx)
		assert.Equal(t, "0102030405060708090a0b0c0d0e0f10", traceID)
	})

	t.Run("case 3", func(t *testing.T) {
		var ctx context.Context
		traceID := ExtractTraceID(ctx)
		assert.Equal(t, "", traceID)
	})
}

func TestCustomTag(t *testing.T) {
	out := CustomTag("hello", "")
	in := attribute.KeyValue{
		Key:   "hello",
		Value: attribute.Value{},
	}
	assert.Equal(t, in.Key, out.Key)
}

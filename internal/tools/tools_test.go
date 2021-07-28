package tools

import (
	"context"
	"testing"

	"github.com/gotomicro/ego/core/transport"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestGetContextValue(t *testing.T) {
	md := metadata.New(map[string]string{
		"X-Ego-Uid": "9527",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	value := GrpcHeaderValue(ctx, "X-Ego-Uid")
	assert.Equal(t, "9527", value)
}

func TestLoggerContextValueByHeader(t *testing.T) {
	md := metadata.New(map[string]string{
		"X-Ego-Uid": "9527",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	value := LoggerGrpcContextValue(ctx, "X-Ego-Uid")
	assert.Equal(t, "9527", value)
}

func TestLoggerContextValueByCtxValue(t *testing.T) {
	transport.Set([]string{"X-Ego-Uid"})

	ctx := transport.WithValue(context.Background(), "X-Ego-Uid", 9527)
	value := LoggerGrpcContextValue(ctx, "X-Ego-Uid")
	assert.Equal(t, "9527", value)

	ctx = transport.WithValue(context.Background(), "X-Ego-Uid", 9528)
	value = LoggerGrpcContextValue(ctx, "X-Ego-Uid")
	assert.Equal(t, "9528", value)
}

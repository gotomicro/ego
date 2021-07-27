package tools

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestGetContextValue(t *testing.T) {
	md := metadata.New(map[string]string{
		"X-Ego-Uid": "9527",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	value := GetContextValue(ctx, "X-Ego-Uid")
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
	ctx := context.WithValue(context.Background(), "X-Ego-Uid", 9527)
	value := LoggerGrpcContextValue(ctx, "X-Ego-Uid")
	assert.Equal(t, "9527", value)
}

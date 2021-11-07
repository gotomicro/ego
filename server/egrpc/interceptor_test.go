package egrpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func Test_getPeerName(t *testing.T) {
	md := metadata.New(map[string]string{
		"app": "ego-svc",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	value := getPeerName(ctx)
	assert.Equal(t, "ego-svc", value)

	ctx2 := metadata.NewIncomingContext(context.Background(), nil)
	value2 := getPeerName(ctx2)
	assert.Equal(t, "", value2)
}

// todo add more unittest
func Test_getPeerIP(t *testing.T) {
	md := metadata.New(map[string]string{
		"client-ip": "127.0.0.1",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	value := getPeerIP(ctx)
	assert.Equal(t, "127.0.0.1", value)
}

func Test_enableCPUUsage(t *testing.T) {
	md := metadata.New(map[string]string{
		"enable-cpu-usage": "true",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	value := enableCPUUsage(ctx)
	assert.Equal(t, true, value)

	ctx2 := metadata.NewIncomingContext(context.Background(), nil)
	value2 := enableCPUUsage(ctx2)
	assert.Equal(t, false, value2)

	md3 := metadata.New(map[string]string{
		"enable-cpu-usage": "test",
	})
	ctx3 := metadata.NewIncomingContext(context.Background(), md3)
	value3 := enableCPUUsage(ctx3)
	assert.Equal(t, false, value3)
}

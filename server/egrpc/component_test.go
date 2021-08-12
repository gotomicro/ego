package egrpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/elog"
)

func TestNewComponent(t *testing.T) {
	cfg := Config{
		Host:    "0.0.0.0",
		Port:    9005,
		Network: "tcp4",
	}
	cmp := newComponent("test-cmp", &cfg, elog.DefaultLogger)
	assert.Equal(t, "test-cmp", cmp.Name())
	assert.Equal(t, "server.egrpc", cmp.PackageName())
	assert.Equal(t, "0.0.0.0:9005", cmp.Address())

	assert.NoError(t, cmp.Init())

	info := cmp.Info()
	assert.NotEmpty(t, info.Name)
	assert.Equal(t, "grpc", info.Scheme)
	assert.Equal(t, "0.0.0.0:9005", info.Address)
	assert.Equal(t, constant.ServiceProvider, info.Kind)

	// err = cmp.Start()
	go func() {
		assert.NoError(t, cmp.Start())
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	<-ctx.Done()
	assert.NoError(t, cmp.Stop())

	t.Log("done")
}

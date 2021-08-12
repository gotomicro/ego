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
	logger := elog.DefaultLogger
	cmp := newComponent("test-cmp", &cfg, logger)
	name := cmp.Name()
	assert.Equal(t, "test-cmp", name)

	pkgName := cmp.PackageName()
	assert.Equal(t, "server.egrpc", pkgName)

	addr := cmp.Address()
	assert.Equal(t, "0.0.0.0:9005", addr)

	var err error
	err = cmp.Init()
	assert.NoError(t, err)

	info := cmp.Info()
	assert.NotEmpty(t, info.Name)
	assert.Equal(t, "grpc", info.Scheme)
	assert.Equal(t, "0.0.0.0:9005", info.Address)
	assert.Equal(t, constant.ServiceProvider, info.Kind)

	// err = cmp.Start()
	go func() {
		err := cmp.Start()
		assert.NoError(t, err)
	}()

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
L:
	for {
		select {
		case <-ctx.Done():
			err := cmp.Stop()
			assert.NoError(t, err)
			break L
		}
	}

	t.Log("done")
}

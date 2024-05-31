package egovernor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/elog"
)

func TestComponent(t *testing.T) {
	cfg := Config{
		Host:    "0.0.0.0",
		Port:    9001,
		Network: "tcp4",
	}
	c := newComponent("test", &cfg, elog.DefaultLogger)
	assert.Equal(t, "test", c.Name())
	assert.Equal(t, PackageName, c.PackageName())
	assert.NoError(t, c.Init())

	info := c.Info()
	assert.NotEmpty(t, info.Name)
	assert.Equal(t, "http", info.Scheme)
	assert.Equal(t, "0.0.0.0:9001", info.Address)
	assert.Equal(t, constant.ServiceGovernor, info.Kind)

	go func() {
		assert.NoError(t, c.Start())
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	<-ctx.Done()
	assert.NoError(t, c.Stop())
	assert.NoError(t, c.GracefulStop(context.Background()))

	t.Log("done")
}

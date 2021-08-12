package egrpc

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/econf"
)

func TestDefaultContainer(t *testing.T) {
	c := DefaultContainer()
	assert.NotPanics(t, func() {
		cmp := c.Build()
		addr := cmp.Address()
		assert.Equal(t, ":9002", addr)
	})
}

func TestNewContainer(t *testing.T) {
	cfg := `
[grpc]
port = 9005
host = "127.0.0.1"
`
	err := econf.LoadFromReader(strings.NewReader(cfg), toml.Unmarshal)
	assert.NoError(t, err)
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	assert.NotPanics(t, func() {
		cmp := Load("grpc").Build()
		addr := cmp.Address()
		assert.Equal(t, "127.0.0.1:9005", addr)
		assert.NoError(t, cmp.Init())
		go func() {
			assert.NoError(t, cmp.Start())
		}()
	L:
		for {
			select {
			case <-ctx.Done():
				assert.NoError(t, cmp.Stop())
				break L
			}
		}
	})
	t.Log("done")
}

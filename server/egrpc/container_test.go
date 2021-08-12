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
		assert.Equal(t, ":9002", cmp.Address())
	})
}

func TestNewContainer(t *testing.T) {
	cfg := `
[grpc]
port = 9005
host = "127.0.0.1"
`
	assert.NoError(t, econf.LoadFromReader(strings.NewReader(cfg), toml.Unmarshal))
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	assert.NotPanics(t, func() {
		cmp := Load("grpc").Build()
		assert.Equal(t, "127.0.0.1:9005", cmp.Address())
		assert.NoError(t, cmp.Init())
		go func() {
			assert.NoError(t, cmp.Start())
		}()
		<-ctx.Done()
		assert.NoError(t, cmp.Stop())
	})
	t.Log("done")
}

package egrpc

import (
	"context"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/gotomicro/ego/core/elog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/gotomicro/ego/core/econf"
)

func newCmp(t *testing.T, opt Option) *Component {
	cfg := `
[grpc]
port = 9005
host = "127.0.0.1"
`
	err := econf.LoadFromReader(strings.NewReader(cfg), toml.Unmarshal)
	assert.NoError(t, err)
	cmp := Load("grpc").Build(opt)
	return cmp
}

func TestWithServerOption(t *testing.T) {
	opt := WithServerOption(grpc.WriteBufferSize(128 * 1024))
	cmp := newCmp(t, opt)
	assert.Equal(t, 3, len(cmp.config.serverOptions))
}

func TestWithStreamInterceptor(t *testing.T) {
	intcp := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return nil
	}
	opt := WithStreamInterceptor(intcp)
	cmp := newCmp(t, opt)
	assert.Equal(t, 2, len(cmp.config.streamInterceptors))
	t.Log("done")
}

func TestWithUnaryInterceptor(t *testing.T) {
	intcp := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		return nil, nil
	}
	opt := WithUnaryInterceptor(intcp)
	cmp := newCmp(t, opt)
	assert.Equal(t, 2, len(cmp.config.unaryInterceptors))
	t.Log("done")
}

func TestWithNetwork(t *testing.T) {
	cmp := newCmp(t, WithNetwork("bufnet"))
	assert.Equal(t, "bufnet", cmp.config.Network)
}

func TestWithLogger(t *testing.T) {
	logger := elog.DefaultContainer().Build(
		elog.WithDebug(false),
		elog.WithEnableAddCaller(true),
		elog.WithEnableAsync(false),
	)

	comp := DefaultContainer().Build(WithLogger(logger))
	assert.Equal(t, logger, comp.logger)
}

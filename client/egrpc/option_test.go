package egrpc

import (
	"strings"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gotomicro/ego/core/econf"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/test/bufconn"
)

func newCmp(t *testing.T, opt Option) *Component {
	cfg := `
[grpc]
network="bufnet"
`
	err := econf.LoadFromReader(strings.NewReader(cfg), toml.Unmarshal)
	assert.NoError(t, err)
	// 必须使用bufnet，这样才能启动起来
	cmp := Load("grpc").Build(WithBufnetServerListener(svc.Listener().(*bufconn.Listener)), opt)
	return cmp
}

func TestWithAddr(t *testing.T) {
	opt := WithAddr("127.0.0.1")
	cmp := newCmp(t, opt)
	assert.Equal(t, "127.0.0.1", cmp.config.Addr)
}

func TestWithBalancerName(t *testing.T) {
	opt := WithBalancerName("round_robin")
	cmp := newCmp(t, opt)
	assert.Equal(t, "round_robin", cmp.config.BalancerName)
}

func TestWithDebug(t *testing.T) {
	_ = WithDebug(true)
}

func TestWithEnableAccessInterceptor(t *testing.T) {
	opt := WithEnableAccessInterceptor(true)
	cmp := newCmp(t, opt)
	assert.Equal(t, true, cmp.config.EnableAccessInterceptor)
}

func TestWithEnableAccessInterceptorReq(t *testing.T) {
	opt := WithEnableAccessInterceptorReq(true)
	cmp := newCmp(t, opt)
	assert.Equal(t, true, cmp.config.EnableAccessInterceptorReq)
}

func TestWithEnableAccessInterceptorRes(t *testing.T) {
	opt := WithEnableAccessInterceptorRes(true)
	cmp := newCmp(t, opt)
	assert.Equal(t, true, cmp.config.EnableAccessInterceptorRes)
}

func TestWithOnFail(t *testing.T) {
	opt := WithOnFail("error")
	cmp := newCmp(t, opt)
	assert.Equal(t, "error", cmp.config.OnFail)
}

func TestWithReadTimeout(t *testing.T) {
	opt := WithReadTimeout(1 * time.Second)
	cmp := newCmp(t, opt)
	assert.Equal(t, 1*time.Second, cmp.config.ReadTimeout)
}

package egin

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/elog"
	"github.com/stretchr/testify/assert"
)

func TestInterceptor(t *testing.T) {
	comp := DefaultContainer().Build()
	// healthcheck，默认中间件，监控中间件，限流中间件
	assert.Equal(t, 3, len(comp.Handlers))
}

func TestWithTrustedPlatform(t *testing.T) {
	comp := DefaultContainer().Build(WithTrustedPlatform("X-Custom-IP"))
	assert.Equal(t, "X-Custom-IP", comp.config.TrustedPlatform)
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

func TestWithHost(t *testing.T) {
	comp := DefaultContainer().Build(WithHost("192.168.10.1"))
	assert.Equal(t, "192.168.10.1", comp.config.Host)
}

func TestWithPort(t *testing.T) {
	comp := DefaultContainer().Build(WithPort(8080))
	assert.Equal(t, 8080, comp.config.Port)
}

func TestWithNetwork(t *testing.T) {
	comp := DefaultContainer().Build(WithNetwork("tcp"))
	assert.Equal(t, "tcp", comp.config.Network)
}

func TestWithSentinelResourceExtractor(t *testing.T) {
	comp := DefaultContainer().Build(WithSentinelResourceExtractor(func(c *gin.Context) string {
		return "test"
	}))
	assert.Equal(t, "test", comp.config.resourceExtract(&gin.Context{}))
}

func TestWithTLSSessionCache(t *testing.T) {
	comp := DefaultContainer().Build(WithTLSSessionCache(nil))
	assert.Nil(t, comp.config.TLSSessionCache)
}

func TestWithTrustedPlatform2(t *testing.T) {
	comp := DefaultContainer().Build(WithTrustedPlatform("X-Custom-IP"))
	assert.Equal(t, "X-Custom-IP", comp.config.TrustedPlatform)
}

func TestWithLogger2(t *testing.T) {
	logger := elog.DefaultContainer().Build(
		elog.WithDebug(false),
		elog.WithEnableAddCaller(true),
		elog.WithEnableAsync(false),
	)

	comp := DefaultContainer().Build(WithLogger(logger))
	assert.Equal(t, logger, comp.logger)

}

func TestWithServerReadTimeout(t *testing.T) {
	timeout := time.Duration(100)
	comp := DefaultContainer().Build(WithServerReadTimeout(timeout))
	assert.Equal(t, timeout, comp.config.ServerReadTimeout)
}

func TestWithServerWriteTimeout(t *testing.T) {
	timeout := time.Duration(100)
	comp := DefaultContainer().Build(WithServerWriteTimeout(timeout))
	assert.Equal(t, timeout, comp.config.ServerWriteTimeout)
}

func TestWithContextTimeout(t *testing.T) {
	timeout := time.Duration(100)
	comp := DefaultContainer().Build(WithContextTimeout(timeout))
	assert.Equal(t, timeout, comp.config.ContextTimeout)
}

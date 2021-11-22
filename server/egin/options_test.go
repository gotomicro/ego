package egin

import (
	"testing"

	"github.com/gotomicro/ego/core/elog"
	"github.com/stretchr/testify/assert"
)

func TestInterceptor(t *testing.T) {
	comp := DefaultContainer().Build()
	// healthcheck，默认中间件，监控中间件
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

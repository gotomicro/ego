package egin

import (
	"testing"

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

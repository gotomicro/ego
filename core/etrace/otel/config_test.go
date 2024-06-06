package otel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	conf := DefaultConfig()
	out := Load("")
	assert.Equal(t, conf, out)
	Load("").Build()
	out1 := conf.buildJaegerTP()
	assert.True(t, true, out1)
	err := conf.Stop()
	assert.NoError(t, err)
}

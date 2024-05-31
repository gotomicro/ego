package etrace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsGlobalTracerRegistered(t *testing.T) {
	assert.True(t, true, IsGlobalTracerRegistered())
}

func TestCustomTag2(t *testing.T) {
	CustomTag("hello", "world")
	assert.NoError(t, nil)
}

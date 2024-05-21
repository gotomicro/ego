package otel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildOtlpTP(t *testing.T) {
	Load("").Build()
	assert.NoError(t, nil)
	c := DefaultConfig()
	c.buildJaegerTP()
	assert.NoError(t, nil)
	err := c.Stop()
	assert.NoError(t, err)
}

package esentinel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	c := DefaultConfig()
	assert.Equal(t, "./logs", c.LogPath)
	assert.Equal(t, "", c.FlowRulesFile)
}

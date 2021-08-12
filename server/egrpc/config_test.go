package egrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.Equal(t, "tcp4", cfg.Network)
	assert.Equal(t, 9002, cfg.Port)
	assert.Equal(t, ":9002", cfg.Address())
}

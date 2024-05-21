package egovernor

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/eflag"
)

func TestDefaultConfig(t *testing.T) {
	in := &Config{
		Host:    eflag.String("host"),
		Network: "tcp4",
		Port:    9003,
	}
	out := DefaultConfig()
	assert.Equal(t, in, out)
}

func TestAddress(t *testing.T) {
	config := Config{Host: "hello", Port: 111, EnableLocalMainIP: true, Network: "tcp4"}
	out := config.Address()
	assert.Equal(t, "hello:111", out)
}

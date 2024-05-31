package egovernor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithHost(t *testing.T) {
	c := &Container{config: &Config{Host: "test"}}
	opt := WithHost("test")
	opt(c)
	assert.Equal(t, "test", c.config.Host)
}

func TestWithPost(t *testing.T) {
	c := &Container{config: &Config{Port: 8080}}
	opt := WithPort(8080)
	opt(c)
	assert.Equal(t, 8080, c.config.Port)
}

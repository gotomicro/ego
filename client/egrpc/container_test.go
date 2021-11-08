package egrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultContainer(t *testing.T) {
	c := DefaultContainer()
	assert.Panics(t, func() {
		c.Build()
	})
}

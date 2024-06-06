package econf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainer(t *testing.T) {
	assert.Equal(t, "mapstructure", GetOptionTagName())
	assert.Equal(t, false, GetOptionWeaklyTypedInput())
	assert.Equal(t, false, GetOptionSquash())
}

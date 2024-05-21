package econf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainer(t *testing.T) {
	out1 := GetOptionTagName()
	assert.Equal(t, "mapstructure", out1)

	out2 := GetOptionWeaklyTypedInput()
	assert.Equal(t, false, out2)

	out3 := GetOptionSquash()
	assert.Equal(t, false, out3)
}

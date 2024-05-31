package transport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var mc = &GrpcHeaderCarrier{}

func TestGet(t *testing.T) {
	key := "test"
	assert.Equal(t, "", mc.Get(key))
}

func TestSet(t *testing.T) {
	mc.Set("hello", "world")
	assert.Nil(t, nil)
}

func TestKeys(t *testing.T) {
	out := mc.Keys()
	assert.Equal(t, []string{"hello"}, out)
}

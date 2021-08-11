package xstring

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateUUID(t *testing.T) {
	out := GenerateUUID(time.Now())
	assert.NotEmpty(t, out)
	assert.Equal(t, 32, len(out))
	assert.NotEqual(t, "00000000000000000000000000000000", out)
}

func TestGenerateID(t *testing.T) {
	out := GenerateID()
	assert.NotEmpty(t, out)
	assert.Equal(t, 16, len(out))
	assert.NotEqual(t, "00000000000000000000000000000000", out)
}

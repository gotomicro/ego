package egin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad2(t *testing.T) {
	Load("").Build()
	assert.NoError(t, nil)
}

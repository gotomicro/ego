package egrpclog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	Build()
	assert.NoError(t, nil)
}

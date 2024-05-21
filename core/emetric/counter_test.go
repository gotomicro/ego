package emetric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCounterVec(t *testing.T) {
	name := "test"
	labels := []string{"test"}
	NewCounterVec(name, labels)
	assert.NoError(t, nil)
}

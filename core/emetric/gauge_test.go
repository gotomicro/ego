package emetric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGaugeVec(t *testing.T) {
	name := "test_"
	labels := []string{"hello_", "world_"}
	NewGaugeVec(name, labels)
	assert.NoError(t, nil)
}

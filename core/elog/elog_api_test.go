package elog

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestElogAPI(t *testing.T) {
	f := zap.Field{
		Key:       "test",
		Type:      9,
		Integer:   11,
		String:    "test",
		Interface: nil,
	}
	Info("", f)
	assert.NoError(t, nil)

	Debug("", f)
	assert.NoError(t, nil)

	Warn("", f)
	assert.NoError(t, nil)

	Error("", f)
	assert.NoError(t, nil)

}

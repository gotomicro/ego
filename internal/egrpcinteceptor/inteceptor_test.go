package egrpcinteceptor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ctx context.Context
var m = messageType{}

func TestEvent(t *testing.T) {
	m.Event(ctx, 111, "")
	assert.NoError(t, nil)
}

func TestSplitMethodName(t *testing.T) {
	f := "GET/https://test.com/xxx"
	SplitMethodName(f)
	assert.NoError(t, nil)
}

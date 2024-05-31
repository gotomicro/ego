package transport

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomContextKeys(t *testing.T) {
	Set([]string{"X-EGO-Test"})
	arr := CustomContextKeys()
	assert.Equal(t, []string{"X-EGO-Test"}, arr)
	length := CustomContextKeysLength()
	assert.Equal(t, 1, length)
}

func TestValue(t *testing.T) {
	Set([]string{"X-EGO-Test"})
	ctx := context.Background()
	ctx = WithValue(ctx, "X-EGO-Test", "hello")
	val := ctx.Value("X-EGO-Test")
	assert.Equal(t, "hello", val)
	Value(ctx, "test")
	assert.NoError(t, nil)
}

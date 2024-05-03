package transport

import (
	"context"

	"github.com/gotomicro/ego/core/eapp"
)

var customKeyStore = contextKeyStore{
	keyArr: make([]string, 0),
}

type contextKeyStore struct {
	keyArr []string
	length int
}

func init() {
	customKeyStore.keyArr = eapp.EgoLogExtraKeys()
	customKeyStore.length = len(customKeyStore.keyArr)
}

// Set overrides custom keys with provided array.
func Set(arr []string) {
	length := len(arr)
	customKeyStore.keyArr = arr
	customKeyStore.length = length
}

// CustomContextKeys returns custom content key list
func CustomContextKeys() []string {
	return customKeyStore.keyArr
}

// CustomContextKeysLength returns custom content key list length
func CustomContextKeysLength() int {
	return customKeyStore.length
}

// WithValue returns a new context with your key and value
func WithValue(ctx context.Context, key string, value interface{}) context.Context {
	info := ctx.Value(key)
	if info != nil {
		return ctx
	}
	return context.WithValue(ctx, key, value)
}

// Value returns value of your key
// Deprecated
// Use ctx.Value()
func Value(ctx context.Context, key string) interface{} {
	return ctx.Value(key)
}

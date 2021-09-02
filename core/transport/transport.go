package transport

import (
	"context"

	"github.com/gotomicro/ego/core/eapp"
)

var customKeyStore = contextKeyStore{
	keyArr: make([]string, 0),
	keyMap: make(map[string]*contextKey),
}

type contextKeyStore struct {
	keyArr []string
	keyMap map[string]*contextKey
	length int
}

func init() {
	customKeyStore.keyArr = eapp.EgoLogExtraKeys()
	for _, value := range eapp.EgoLogExtraKeys() {
		customKeyStore.keyMap[value] = newContextKey(value)
	}
	customKeyStore.length = len(customKeyStore.keyArr)
}

// Set 设置context key arr
func Set(arr []string) {
	customKeyStore.keyArr = arr
	for _, value := range arr {
		customKeyStore.keyMap[value] = newContextKey(value)
	}
	customKeyStore.length = len(customKeyStore.keyArr)
}

// CustomContextKeys 自定义context
func CustomContextKeys() []string {
	return customKeyStore.keyArr
}

// CustomContextKeysLength 自定义context keys长度
func CustomContextKeysLength() int {
	return customKeyStore.length
}

// WithValue 设置数据
func WithValue(ctx context.Context, key string, value interface{}) context.Context {
	return context.WithValue(ctx, getContextKey(key), value)
}

// Value 获取数据
func Value(ctx context.Context, key string) interface{} {
	return ctx.Value(getContextKey(key))
}

func newContextKey(name string) *contextKey {
	return &contextKey{name: name}
}

func getContextKey(key string) *contextKey {
	return customKeyStore.keyMap[key]
}

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "ego context value " + k.name }

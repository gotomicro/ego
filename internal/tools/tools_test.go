package tools

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	"github.com/gotomicro/ego/core/transport"
)

func TestGrpcHeaderValue(t *testing.T) {
	value := GrpcHeaderValue(context.Background(), "")
	assert.Equal(t, "", value)

	md := metadata.New(map[string]string{
		"X-Ego-Uid": "9527",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	value2 := GrpcHeaderValue(ctx, "X-Ego-Uid")
	assert.Equal(t, "9527", value2)
}

func TestGrpcHeaderValueEmpty(t *testing.T) {
	value := GrpcHeaderValue(context.Background(), "X-Ego-Uid")
	assert.Equal(t, "", value)
}

func TestContextValue(t *testing.T) {
	value := ContextValue(context.Background(), "")
	assert.Equal(t, "", value)

	transport.Set([]string{"X-Ego-Uid"})
	ctx := transport.WithValue(context.Background(), "X-Ego-Uid", 9527)
	value = ContextValue(ctx, "X-Ego-Uid")
	assert.Equal(t, "9527", value)

	ctx = transport.WithValue(context.Background(), "X-Ego-Uid", 9528)
	value = ContextValue(ctx, "X-Ego-Uid")
	assert.Equal(t, "9528", value)
}

func TestContextValueEmpty(t *testing.T) {
	value := ContextValue(context.Background(), "X-Ego-Uid")
	assert.Equal(t, "", value)
}

func TestToSliceStringMap(t *testing.T) {
	out := ToSliceStringMap([]interface{}{
		map[string]interface{}{"aaa": "AAA"},
	})
	assert.Equal(t, []map[string]interface{}{{"aaa": "AAA"}}, out)
}

func TestGofmt(t *testing.T) {
	assert.Panics(t, func() {
		GoFmt([]byte(`asdfasdfasdfasdfasd func main`))
	})
}

package transport

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestGetAndKeys(t *testing.T) {
	md := metadata.Pairs("hello", "world", "bye", "test")
	carrier := GrpcHeaderCarrier(md)

	// 测试 Get()
	t.Run("case 1", func(t *testing.T) {
		out := carrier.Get("hello")
		assert.Equal(t, "world", out)
	})

	t.Run("case 2", func(t *testing.T) {
		out := carrier.Get("testing")
		assert.Equal(t, "", out)
	})

	// 测试Keys()
	keys := carrier.Keys()
	reflect.DeepEqual([]string{"hello", "bye"}, keys)
}

func TestSet(t *testing.T) {
	md := metadata.MD{}
	carrier := GrpcHeaderCarrier(md)
	carrier.Set("hello", "world")
	out := carrier.Get("hello")
	assert.Equal(t, "world", out)
}

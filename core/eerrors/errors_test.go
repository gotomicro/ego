package eerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestRegister(t *testing.T) {
	errUnknown := New(int(codes.Unknown), "unknown", "unknown")
	Register(errUnknown)

	// 一个新error，添加信息
	newErrUnknown := errUnknown.WithMessage("unknown something").WithMetadata(map[string]string{
		"hello": "world",
	}).(*EgoError)
	assert.Equal(t, "unknown something", newErrUnknown.GetMessage())
	assert.Equal(t, map[string]string{
		"hello": "world",
	}, newErrUnknown.GetMetadata())

	assert.ErrorIs(t, newErrUnknown, errUnknown)
}

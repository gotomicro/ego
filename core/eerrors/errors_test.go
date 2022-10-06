package eerrors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

var (
	nilEgoErr    *EgoError
	nilErr       error
	notNilEgoErr = New(1, "__REASON__", "__MESSAGE__")
	errNotNil    = errors.New("some error")
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

func TestIs(t *testing.T) {
	tests := []struct {
		name        string
		originalErr *EgoError
		targetErr   error
		wantRes     bool
	}{
		{"nilEgoErr-nilEgoErr", nilEgoErr, nilEgoErr, true},
		{"nilEgoErr-nilErr", nilEgoErr, nilErr, false},
		{"notNilEgoErr-errNotNil", nilEgoErr, errNotNil, false},
		{"notNilEgoErr-notNilEgoErr", nilEgoErr, notNilEgoErr, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantRes, tt.originalErr.Is(tt.targetErr))
			assert.Equal(t, tt.wantRes, errors.Is(tt.originalErr, tt.targetErr))
		})
	}
}

func TestFromError(t *testing.T) {
	assert.Equal(t, nilEgoErr, FromError(nil))
	tests := []struct {
		name         string
		req          error
		wantCode     int32
		wantReason   string
		wantMessage  string
		wantMetadata map[string]string
	}{
		{"empty-ego-error", New(0, "", ""), 0, "", "", nil},
		{"some-ego-error", notNilEgoErr, 1, "__REASON__", "__MESSAGE__", nil},
		{"normal-error", errNotNil, int32(codes.Unknown), UnknownReason, "some error", nil},
	}
	for _, tt := range tests {
		res := FromError(tt.req)
		assert.Equal(t, tt.wantCode, res.Code)
		assert.Equal(t, tt.wantReason, res.Reason)
		assert.Equal(t, tt.wantMessage, res.Message)
		assert.Equal(t, tt.wantMetadata, res.Metadata)
	}
}

func TestCheckErr(t *testing.T) {
	tests := []struct {
		name        string
		originalErr error
		checkedErr  *EgoError
		wantRes     bool
	}{
		{"original:nil, checked:nil", nil, nil, true},
		{"original:notNilEgoErr, checked:nil", notNilEgoErr, nil, false},
		{"original:errNotNil, checked:nil", errNotNil, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantRes, errors.Is(FromError(tt.originalErr), tt.checkedErr))
		})
	}
}

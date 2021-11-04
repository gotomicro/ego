package ecode

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestConvert(t *testing.T) {
	assert.Equal(t, status.New(codes.OK, "OK"), Convert(nil))
	assert.Equal(t, status.New(codes.Canceled, context.Canceled.Error()), Convert(context.Canceled))
	assert.Equal(t, status.New(codes.DeadlineExceeded, context.DeadlineExceeded.Error()), Convert(context.DeadlineExceeded))
}

func TestGrpcToHTTPStatusCode(t *testing.T) {
	code := GrpcToHTTPStatusCode(status.New(codes.OK, "OK").Code())
	assert.Equal(t, http.StatusOK, code)
}

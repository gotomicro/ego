package ecode

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestConvert(t *testing.T) {
	statusInfo := Convert(nil)
	assert.Equal(t, status.New(codes.OK, "OK"), statusInfo)
}

func TestGrpcToHTTPStatusCode(t *testing.T) {
	code := GrpcToHTTPStatusCode(status.New(codes.OK, "OK").Code())
	assert.Equal(t, http.StatusOK, code)
}

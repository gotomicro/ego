package ecode

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type errTest struct {
}

func (*errTest) Error() string {
	return "i am error"
}

func (*errTest) GRPCStatus() *status.Status {
	return status.New(codes.Code(999), "i am error grpc status")
}

func TestConvert(t *testing.T) {
	errtt := &errTest{}
	assert.Equal(t, status.New(codes.Code(999), "i am error grpc status"), Convert(errtt))
	assert.Equal(t, status.New(codes.OK, "OK"), Convert(nil))
	assert.Equal(t, status.New(codes.Canceled, context.Canceled.Error()), Convert(context.Canceled))
	assert.Equal(t, status.New(codes.DeadlineExceeded, context.DeadlineExceeded.Error()), Convert(context.DeadlineExceeded))
	unknownErr := fmt.Errorf("unknown error")
	assert.Equal(t, status.New(codes.Unknown, unknownErr.Error()), Convert(unknownErr))
}

func TestGrpcToHTTPStatusCode(t *testing.T) {
	type args struct {
		statusCode codes.Code
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				statusCode: status.New(codes.OK, "").Code(),
			},
			want: http.StatusOK,
		},
		{
			args: args{
				statusCode: status.New(codes.Canceled, "").Code(),
			},
			want: http.StatusRequestTimeout,
		},
		{
			args: args{
				statusCode: status.New(codes.Unknown, "").Code(),
			},
			want: http.StatusInternalServerError,
		},
		{
			args: args{
				statusCode: status.New(codes.InvalidArgument, "").Code(),
			},
			want: http.StatusBadRequest,
		},
		{
			args: args{
				statusCode: status.New(codes.DeadlineExceeded, "").Code(),
			},
			want: http.StatusRequestTimeout,
		},
		{
			args: args{
				statusCode: status.New(codes.NotFound, "").Code(),
			},
			want: http.StatusNotFound,
		},
		{
			args: args{
				statusCode: status.New(codes.AlreadyExists, "").Code(),
			},
			want: http.StatusConflict,
		},
		{
			args: args{
				statusCode: status.New(codes.PermissionDenied, "").Code(),
			},
			want: http.StatusForbidden,
		},
		{
			args: args{
				statusCode: status.New(codes.Unauthenticated, "").Code(),
			},
			want: http.StatusUnauthorized,
		},
		{
			args: args{
				statusCode: status.New(codes.ResourceExhausted, "").Code(),
			},
			want: http.StatusServiceUnavailable,
		},
		{
			args: args{
				statusCode: status.New(codes.FailedPrecondition, "").Code(),
			},
			want: http.StatusPreconditionFailed,
		},
		{
			args: args{
				statusCode: status.New(codes.Aborted, "").Code(),
			},
			want: http.StatusConflict,
		},
		{
			args: args{
				statusCode: status.New(codes.OutOfRange, "").Code(),
			},
			want: http.StatusBadRequest,
		},
		{
			args: args{
				statusCode: status.New(codes.Unimplemented, "").Code(),
			},
			want: http.StatusNotImplemented,
		},
		{
			args: args{
				statusCode: status.New(codes.Internal, "").Code(),
			},
			want: http.StatusInternalServerError,
		},
		{
			args: args{
				statusCode: status.New(codes.Unavailable, "").Code(),
			},
			want: http.StatusServiceUnavailable,
		},
		{
			args: args{
				statusCode: status.New(codes.DataLoss, "").Code(),
			},
			want: http.StatusInternalServerError,
		},
		{
			args: args{
				statusCode: status.New(codes.Code(9999), "").Code(),
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GrpcToHTTPStatusCode(tt.args.statusCode); got != tt.want {
				t.Errorf("GrpcToHTTPStatusCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

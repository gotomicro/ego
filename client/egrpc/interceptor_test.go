package egrpc

import (
	"context"
	"testing"

	"github.com/gotomicro/ego/internal/tools"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Test_customHeader(t *testing.T) {
	md := metadata.New(map[string]string{
		"X-Ego-Uid": "9527",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	interceptor := customHeader([]string{"X-Ego-Uid"})

	cc := new(grpc.ClientConn)
	err := interceptor(ctx, "/foo", nil, nil, cc,
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			info := tools.GetContextValue(ctx, "X-Ego-Uid")
			assert.Equal(t, "9527", info)
			return nil
		})
	assert.Nil(t, err)
}

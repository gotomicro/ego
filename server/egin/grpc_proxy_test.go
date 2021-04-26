package egin

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

type GreeterMock struct{}

func (mock GreeterMock) SayHello(context context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{
		Message: "hello",
	}, nil
}

func TestGRPCProxyWrapper(t *testing.T) {
	router := gin.New()
	mock := GreeterMock{}
	router.POST("/", GRPCProxy(mock.SayHello))

	// RUN
	w := performRequest(router, "POST", "/")

	assert.Equal(t, 200, w.Code)
}

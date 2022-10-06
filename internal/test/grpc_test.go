package test

import (
	"context"
	"log"
	"testing"

	cegrpc "github.com/gotomicro/ego/client/egrpc"
	"github.com/gotomicro/ego/core/eerrors"
	"github.com/gotomicro/ego/internal/test/errcode"
	"github.com/gotomicro/ego/internal/test/helloworld"
	"github.com/gotomicro/ego/server/egrpc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
)

var svc *egrpc.Component

func init() {
	svc = egrpc.DefaultContainer().Build(egrpc.WithNetwork("bufnet"))
	helloworld.RegisterGreeterServer(svc, &Greeter{})
	err := svc.Init()
	if err != nil {
		log.Fatalf("init exited with error: %v", err)
	}
	go func() {
		err = svc.Start()
		if err != nil {
			log.Fatalf("init start with error: %v", err)
		}
	}()
}

func TestGrpcError(t *testing.T) {
	resourceClient := cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener().(*bufconn.Listener)))
	ctx := context.Background()
	client := helloworld.NewGreeterClient(resourceClient.ClientConn)
	_, err := client.SayHello(ctx, &helloworld.HelloRequest{})
	egoErr := eerrors.FromError(err)
	assert.ErrorIs(t, egoErr, errcode.ErrInvalidArgument())
	assert.Equal(t, "name is empty", egoErr.GetMessage())
}

func TestGrpcOk(t *testing.T) {
	resourceClient := cegrpc.DefaultContainer().Build(cegrpc.WithBufnetServerListener(svc.Listener().(*bufconn.Listener)))
	ctx := context.Background()
	client := helloworld.NewGreeterClient(resourceClient.ClientConn)
	resp, err := client.SayHello(ctx, &helloworld.HelloRequest{
		Name: "Ego",
	})
	assert.NoError(t, err)
	assert.True(t, proto.Equal(&helloworld.HelloResponse{
		Message: "Hello Ego",
	}, resp))

}

// Greeter ...
type Greeter struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello ...
func (g Greeter) SayHello(context context.Context, request *helloworld.HelloRequest) (*helloworld.HelloResponse, error) {
	if request.Name == "" {
		return nil, errcode.ErrInvalidArgument().WithMessage("name is empty")
	}

	return &helloworld.HelloResponse{
		Message: "Hello " + request.Name,
	}, nil
}

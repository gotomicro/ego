package egrpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/internal/test/errcode"
	"github.com/gotomicro/ego/internal/test/helloworld"
	"github.com/gotomicro/ego/server/egrpc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
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

func TestComponent_Error(t *testing.T) {
	c := &Component{
		err: fmt.Errorf("some error"),
	}
	assert.EqualError(t, c.Error(), "some error")
}

func Test_newComponent(t *testing.T) {
	// address为空的时候会panic
	assert.Panics(t, func() {
		cfg := DefaultConfig()
		// cfg.OnFail = "error"
		newComponent("test-cmp", cfg, elog.DefaultLogger)
	})

	cfg := DefaultConfig()
	cfg.dialOptions = append(cfg.dialOptions, grpc.WithContextDialer(bufDialer))
	cmp := newComponent("test-cmp", cfg, elog.DefaultLogger)
	ctx := context.Background()
	client := helloworld.NewGreeterClient(cmp.ClientConn)
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

func bufDialer(context.Context, string) (net.Conn, error) {
	return svc.Listener().(*bufconn.Listener).Dial()
}

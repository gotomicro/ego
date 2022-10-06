package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/gotomicro/ego/core/util/xtime"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gotomicro/ego/examples/helloworld"
)

//  export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
	serverOptions := []grpc.ServerOption{grpc.ChainUnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		fmt.Printf("resp--------------->"+"%+v\n", resp)
		fmt.Printf("err--------------->"+"%+v\n", err)
		return
	}), grpc.StatsHandler(&ocgrpc.ServerHandler{})}

	newServer := grpc.NewServer(serverOptions...)

	helloworld.RegisterGreeterServer(newServer, &Greeter{})
	listener, _ := net.Listen("tcp", ":9002")
	newServer.Serve(listener)
}

// Greeter ...
type Greeter struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello ...
func (g Greeter) SayHello(ctx context.Context, request *helloworld.HelloRequest) (*helloworld.HelloResponse, error) {
	if request.Name == "error" {
		return nil, status.Error(codes.Unavailable, "error")
	}

	time.Sleep(xtime.Duration("2s"))
	return &helloworld.HelloResponse{
		Message: "Hello EGO ",
	}, nil
}

package main

import (
	"context"

	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/examples/helloworld"
	"github.com/gotomicro/ego/server"
	"github.com/gotomicro/ego/server/egrpc"
)

//  export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
	if err := ego.New().Serve(func() server.Server {
		server := egrpc.Load("server.grpc").Build()
		helloworld.RegisterGreeterServer(server.Server, &Greeter{server: server})
		return server
	}()).Run(); err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

// Greeter ...
type Greeter struct {
	server *egrpc.Component
	helloworld.UnimplementedGreeterServer
}

// SayHello ...
func (g Greeter) SayHello(context context.Context, request *helloworld.HelloRequest) (*helloworld.HelloResponse, error) {
	return &helloworld.HelloResponse{
		Message: "Hello EGO, I'm " + g.server.Address(),
	}, nil
}

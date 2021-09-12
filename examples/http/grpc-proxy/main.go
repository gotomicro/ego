package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/examples/helloworld"
	"github.com/gotomicro/ego/server/egin"
)

//  export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
	if err := ego.New().Serve(func() *egin.Component {
		server := egin.Load("server.http").Build()
		server.GET("/hello", func(ctx *gin.Context) {
			ctx.JSON(200, "Hello client: "+ctx.GetHeader("app"))
			return
		})
		mock := &GreeterMock{}
		server.GET("/grpcproxyok", egin.GRPCProxy(mock.SayHelloOK))
		server.GET("/grpcproxyerr", egin.GRPCProxy(mock.SayHelloErr))
		return server
	}()).Run(); err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

type GreeterMock struct{}

func (mock GreeterMock) SayHelloOK(context context.Context, request *helloworld.HelloRequest) (*helloworld.HelloResponse, error) {
	return &helloworld.HelloResponse{
		Message: "hello",
	}, nil
}

func (mock GreeterMock) SayHelloErr(context context.Context, request *helloworld.HelloRequest) (*helloworld.HelloResponse, error) {
	return &helloworld.HelloResponse{
		Message: "hello",
	}, fmt.Errorf("say hello err")
}

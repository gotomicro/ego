package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egin"
	"google.golang.org/grpc/examples/helloworld/helloworld"
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
		server.GET("/grpcproxy", egin.GRPCProxy(mock.SayHello))
		return server
	}()).Run(); err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

type GreeterMock struct{}

func (mock GreeterMock) SayHello(context context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{
		Message: "hello",
	}, fmt.Errorf("error say hello")
}

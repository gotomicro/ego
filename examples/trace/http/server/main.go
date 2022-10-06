package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egin"
)

//  export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
	if err := ego.New().Serve(func() *egin.Component {
		server := egin.Load("server.http").Build()

		server.GET("/panic", func(ctx *gin.Context) {
			<-ctx.Request.Context().Done()
			panic(ctx.Request.Context().Err())
		})

		server.GET("/200", func(ctx *gin.Context) {
			<-ctx.Request.Context().Done()
			fmt.Println(ctx.Request.Context().Err())
			ctx.String(200, "hello")
		})

		server.GET("/hello", func(ctx *gin.Context) {
			ctx.JSON(200, "Hello client: "+ctx.GetHeader("app"))
		})
		return server
	}()).Run(); err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

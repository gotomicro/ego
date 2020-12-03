package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egin"
)

func main() {
	if err := ego.New().
		Serve(func() *egin.Component {
			server := egin.Load("server.http").Build()
			server.GET("/hello", func(ctx *gin.Context) {
				ctx.JSON(200, "Hello")
				return
			})
			return server
		}()).Run(); err != nil {
		elog.Panic("startup", elog.Any("err", err))
	}
}

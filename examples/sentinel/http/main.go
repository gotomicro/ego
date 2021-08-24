package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egin"
	"github.com/gotomicro/ego/server/egovernor"
)

// export EGO_DEBUG=true && go run main.go
// ab -n 10 -c 10  http://127.0.0.1:9007/hello，可以看到429，说明限流
func main() {
	if err := ego.New().Serve(func() *egin.Component {
		server := egin.Load("server.http").Build(
		//egin.WithSentinelResourceExtractor(func(ctx *gin.Context) string {
		//	return ctx.Request.Method + "." + ctx.FullPath()
		//}),
		//egin.WithSentinelBlockFallback(func(ctx *gin.Context) {
		//	ctx.AbortWithStatusJSON(429, gin.H{"msg": "too many requests"})
		//}),
		)
		server.GET("/hello", func(c *gin.Context) {
			c.JSON(200, "Hello EGO")
			return
		})
		return server
	}(),
		egovernor.Load("server.governor").Build(),
	).Run(); err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

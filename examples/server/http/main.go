package main

import (
	"github.com/gin-gonic/gin"

	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egin"
)

// export EGO_DEBUG=true && go run main.go --config=config.toml
// curl -i 'http://localhost:9006/hello?q=query' -X POST -H 'X-Ego-Uid: 9999' --data '{"id":1,"name":"lee"}'
func main() {
	if err := ego.New().Serve(func() *egin.Component {
		server := egin.Load("server.http").Build()
		server.GET("/hello", func(c *gin.Context) {
			c.JSON(200, "Hello EGO")
			return
		})
		server.POST("/hello", func(c *gin.Context) {
			var user struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}
			if err := c.BindJSON(&user); err != nil {
				c.JSON(401, "invalid params")
				return
			}

			c.JSON(200, gin.H{
				"q":       c.Query("q"),
				"xEgoUid": c.GetHeader("x-ego-uid"),
				"user":    user,
			})
			return
		})
		return server
	}()).Run(); err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

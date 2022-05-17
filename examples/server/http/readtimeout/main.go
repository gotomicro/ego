package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/client/ehttp"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egin"
)

// export EGO_DEBUG=true && go run main.go --config=config.toml
// curl -i 'http://localhost:9006/hello?q=query' -X POST -H 'X-Ego-Uid: 9999' --data '{"id":1,"name":"lee"}'
func main() {
	if err := ego.New().Invoker(func() error {
		go func() {
			time.Sleep(1 * time.Second)
			startTime := time.Now()
			ehttp.DefaultContainer().Build().R().Get("http://127.0.0.1:12345/hello")
			fmt.Println("cost: ", time.Now().Sub(startTime))
		}()

		return nil
	}).Serve(func() *egin.Component {
		server := egin.Load("server.http").Build()
		server.GET("/hello", func(c *gin.Context) {
			time.Sleep(10 * time.Second)
			//startTime := time.Now()
			//ehttp.DefaultContainer().Build().R().SetContext(c.Request.Context()).Get("http://127.0.0.1:12345/longtime")
			//fmt.Println("cost: ", time.Now().Sub(startTime))
			c.JSON(200, "Hello EGO")
			return
		})

		server.GET("/longtime", func(c *gin.Context) {
			time.Sleep(10 * time.Second)
			c.JSON(200, "Hello longtime")
			return
		})
		return server
	}()).Run(); err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

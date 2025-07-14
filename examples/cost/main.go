package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// XResponseTimer wrap gin reponse writer add start time
type XResponseTimer struct {
	gin.ResponseWriter
	start time.Time
}

func (w *XResponseTimer) WriteHeader(statusCode int) {
	duration := time.Since(w.start)
	lan := strconv.FormatInt(duration.Milliseconds(), 10) + "ms"
	w.Header().Set("X-Response-Time", lan)
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *XResponseTimer) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

// NewXResponseTimer middleware to add X-Response-Time
func NewXResponseTimer(c *gin.Context) {
	blw := &XResponseTimer{ResponseWriter: c.Writer, start: time.Now()}
	c.Writer = blw
	c.Next()
}
func main() {
	r := gin.Default()

	// 全局使用耗时中间件
	r.Use(NewXResponseTimer)

	r.GET("/fast", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "fast response"})
	})

	r.GET("/slow", func(c *gin.Context) {
		time.Sleep(500 * time.Millisecond)
		c.JSON(200, gin.H{"message": "slow response"})
	})

	r.Run(":8082")
}

package egin

import (
	"github.com/gin-gonic/gin"
)

// Option 可选项
type Option func(c *Container)

// WebSocketOption ..
type WebSocketOption func(*WebSocket)

// WithSentinelResourceExtractor 资源命名方式
func WithSentinelResourceExtractor(fn func(*gin.Context) string) Option {
	return func(c *Container) {
		c.config.resourceExtract = fn
	}
}

// WithSentinelBlockFallback 限流后的返回数据
func WithSentinelBlockFallback(fn func(*gin.Context)) Option {
	return func(c *Container) {
		c.config.blockFallback = fn
	}
}

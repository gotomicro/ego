package egin

import (
	"crypto/tls"
	"embed"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/elog"
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

// WithTLSSessionCache 限流后的返回数据
func WithTLSSessionCache(tsc tls.ClientSessionCache) Option {
	return func(c *Container) {
		c.config.TLSSessionCache = tsc
	}
}

// WithTrustedPlatform 信任的Header头，获取客户端IP地址
func WithTrustedPlatform(trustedPlatform string) Option {
	return func(c *Container) {
		c.config.TrustedPlatform = trustedPlatform
	}
}

// WithLogger 信任的Header头，获取客户端IP地址
func WithLogger(logger *elog.Component) Option {
	return func(c *Container) {
		c.logger = logger
	}
}

// WithEmbedFs 设置embed fs
func WithEmbedFs(fs embed.FS) Option {
	return func(c *Container) {
		c.config.embedFs = fs
	}
}

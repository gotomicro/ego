package egin

import (
	"crypto/tls"
	"embed"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gotomicro/ego/core/elog"
)

// Option overrides a Container's default configuration.
type Option func(c *Container)

// WebSocketOption ..
type WebSocketOption func(*WebSocket)

// WithHost 设置host
func WithHost(host string) Option {
	return func(c *Container) {
		c.config.Host = host
	}
}

// WithPort 设置port
func WithPort(port int) Option {
	return func(c *Container) {
		c.config.Port = port
	}
}

// WithNetwork 设置network
func WithNetwork(network string) Option {
	return func(c *Container) {
		c.config.Network = network
	}
}

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

// WithTLSSessionCache TLS Session 缓存
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

// WithLogger 设置 logger
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

// WithServerReadTimeout 设置超时时间
func WithServerReadTimeout(timeout time.Duration) Option {
	return func(c *Container) {
		c.config.ServerReadTimeout = timeout
	}
}

// WithServerReadHeaderTimeout 设置超时时间
func WithServerReadHeaderTimeout(timeout time.Duration) Option {
	return func(c *Container) {
		c.config.ServerReadHeaderTimeout = timeout
	}
}

// WithServerWriteTimeout 设置超时时间
func WithServerWriteTimeout(timeout time.Duration) Option {
	return func(c *Container) {
		c.config.ServerWriteTimeout = timeout
	}
}

// WithContextTimeout 设置 context 超时时间
func WithContextTimeout(timeout time.Duration) Option {
	return func(c *Container) {
		c.config.ContextTimeout = timeout
	}
}

// WithRecoveryFunc 设置 recovery func
func WithRecoveryFunc(f gin.RecoveryFunc) Option {
	return func(c *Container) {
		c.config.recoveryFunc = f
	}
}

func WithListener(listener net.Listener) Option {
	return func(c *Container) {
		c.config.listener = listener
	}
}

func WithCompatibleOtherTrace(f func(http.Header)) Option {
	return func(c *Container) {
		c.config.compatibleTrace = f
	}
}

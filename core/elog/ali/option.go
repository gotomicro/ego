package ali

import (
	"time"

	"go.uber.org/zap/zapcore"
)

// Option 可选项
type Option func(c *config)

// WithEncoder ...
func WithEncoder(enc zapcore.Encoder) Option {
	return func(c *config) {
		c.encoder = enc
	}
}

// WithEndpoint ...
func WithEndpoint(endpoint string) Option {
	return func(c *config) {
		c.endpoint = endpoint
	}
}

// WithAccessKeyID ...
func WithAccessKeyID(akID string) Option {
	return func(c *config) {
		c.accessKeyID = akID
	}
}

// WithAccessKeySecret ....
func WithAccessKeySecret(akSecret string) Option {
	return func(c *config) {
		c.accessKeySecret = akSecret
	}
}

// WithProject ...
func WithProject(project string) Option {
	return func(c *config) {
		c.project = project
	}
}

// WithLogstore ...
func WithLogstore(logStore string) Option {
	return func(c *config) {
		c.logstore = logStore
	}
}

// WithMaxQueueSize ...
func WithMaxQueueSize(maxQueueSize int) Option {
	return func(c *config) {
		c.maxQueueSize = maxQueueSize
	}
}

// WithLevelEnabler ...
func WithLevelEnabler(lv zapcore.LevelEnabler) Option {
	return func(c *config) {
		c.levelEnabler = lv
	}
}

// WithFlushBufferSize ...
func WithFlushBufferSize(flushBufferSize int) Option {
	return func(c *config) {
		c.flushBufferSize = int32(flushBufferSize)
	}
}

// WithFlushBufferInterval ...
func WithFlushBufferInterval(flushBufferInterval time.Duration) Option {
	return func(c *config) {
		c.flushBufferInterval = flushBufferInterval
	}
}

// WithAPIBulkSize ...
func WithAPIBulkSize(apiBulkSize int) Option {
	return func(c *config) {
		c.apiBulkSize = apiBulkSize
	}
}

// WithAPITimeout ...
func WithAPITimeout(apiTimeout time.Duration) Option {
	return func(c *config) {
		c.apiTimeout = apiTimeout
	}
}

// WithAPIRetryCount ...
func WithAPIRetryCount(apiRetryCount int) Option {
	return func(c *config) {
		c.apiRetryCount = apiRetryCount
	}
}

// WithAPIRetryWaitTime ...
func WithAPIRetryWaitTime(apiRetryWaitTime time.Duration) Option {
	return func(c *config) {
		c.apiRetryWaitTime = apiRetryWaitTime
	}
}

// WithAPIRetryMaxWaitTime ...
func WithAPIRetryMaxWaitTime(apiRetryMaxWaitTime time.Duration) Option {
	return func(c *config) {
		c.apiRetryMaxWaitTime = apiRetryMaxWaitTime
	}
}

// WithAPIMaxIdleConns ...
func WithAPIMaxIdleConns(apiMaxIdleConns int) Option {
	return func(c *config) {
		c.apiMaxIdleConns = apiMaxIdleConns
	}
}

// WithAPIIdleConnTimeout ...
func WithAPIIdleConnTimeout(apiIdleConnTimeout time.Duration) Option {
	return func(c *config) {
		c.apiIdleConnTimeout = apiIdleConnTimeout
	}
}

// WithAPIMaxIdleConnsPerHost ...
func WithAPIMaxIdleConnsPerHost(apiMaxIdleConnsPerHost int) Option {
	return func(c *config) {
		c.apiMaxIdleConnsPerHost = apiMaxIdleConnsPerHost
	}
}

// WithFallbackCore ...
func WithFallbackCore(core zapcore.Core) Option {
	return func(c *config) {
		c.fallbackCore = core
	}
}

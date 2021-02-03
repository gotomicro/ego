package ali

import (
	"time"

	"go.uber.org/zap/zapcore"
)

type Option func(c *config)

func WithEncoder(enc zapcore.Encoder) Option {
	return func(c *config) {
		c.encoder = enc
	}
}

func WithEndpoint(endpoint string) Option {
	return func(c *config) {
		c.endpoint = endpoint
	}
}

func WithAccessKeyID(akID string) Option {
	return func(c *config) {
		c.accessKeyID = akID
	}
}

func WithAccessKeySecret(akSecret string) Option {
	return func(c *config) {
		c.accessKeySecret = akSecret
	}
}

func WithProject(project string) Option {
	return func(c *config) {
		c.project = project
	}
}

func WithLogstore(logStore string) Option {
	return func(c *config) {
		c.logstore = logStore
	}
}

func WithLevelEnabler(lv zapcore.LevelEnabler) Option {
	return func(c *config) {
		c.levelEnabler = lv
	}
}

func WithFlushBufferSize(flushBufferSize int) Option {
	return func(c *config) {
		c.flushBufferSize = int32(flushBufferSize)
	}
}

func WithFlushBufferInterval(flushBufferInterval time.Duration) Option {
	return func(c *config) {
		c.flushBufferInterval = flushBufferInterval
	}
}

func WithApiBulkSize(apiBulkSize int) Option {
	return func(c *config) {
		c.apiBulkSize = apiBulkSize
	}
}

func WithApiTimeout(apiTimeout time.Duration) Option {
	return func(c *config) {
		c.apiTimeout = apiTimeout
	}
}

func WithApiRetryCount(apiRetryCount int) Option {
	return func(c *config) {
		c.apiRetryCount = apiRetryCount
	}
}

func WithApiRetryWaitTime(apiRetryWaitTime time.Duration) Option {
	return func(c *config) {
		c.apiRetryWaitTime = apiRetryWaitTime
	}
}

func WithApiRetryMaxWaitTime(apiRetryMaxWaitTime time.Duration) Option {
	return func(c *config) {
		c.apiRetryMaxWaitTime = apiRetryMaxWaitTime
	}
}

func WithFallbackCore(core zapcore.Core) Option {
	return func(c *config) {
		c.fallbackCore = core
	}
}

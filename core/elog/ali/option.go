package ali

import (
	"time"

	"go.uber.org/zap/zapcore"
)

type Option func(c *config)

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

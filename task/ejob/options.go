package ejob

import (
	"context"
)

func WithName(name string) Option {
	return func(c *Container) {
		c.config.Name = name
	}
}

func WithStartFunc(startFunc func(ctx context.Context) error) Option {
	return func(c *Container) {
		c.config.startFunc = startFunc
	}
}

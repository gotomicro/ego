package ejob

func WithName(name string) Option {
	return func(c *Container) {
		c.config.Name = name
	}
}

func WithStartFunc(startFunc func() error) Option {
	return func(c *Container) {
		c.config.startFunc = startFunc
	}
}

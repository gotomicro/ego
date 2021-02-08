package elog

func WithFileName(name string) Option {
	return func(c *Container) {
		c.Config.Name = name
	}
}

func WithDebug(debug bool) Option {
	return func(c *Container) {
		c.Config.Debug = debug
	}
}

func WithLevel(level string) Option {
	return func(c *Container) {
		c.Config.Level = level
	}
}

func WithEnableAsync(enableAsync bool) Option {
	return func(c *Container) {
		c.Config.EnableAsync = enableAsync
	}
}

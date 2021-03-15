package elog

func WithFileName(name string) Option {
	return func(c *Container) {
		c.config.Name = name
	}
}

func WithDebug(debug bool) Option {
	return func(c *Container) {
		c.config.Debug = debug
	}
}

func WithLevel(level string) Option {
	return func(c *Container) {
		c.config.Level = level
	}
}

// WithEnableAsync 是否异步执行，默认异步
func WithEnableAsync(enableAsync bool) Option {
	return func(c *Container) {
		c.config.EnableAsync = enableAsync
	}
}

// WithEnableAddCaller 是否添加行号，默认不添加行号
func WithEnableAddCaller(enableAddCaller bool) Option {
	return func(c *Container) {
		c.config.EnableAddCaller = enableAddCaller
	}
}

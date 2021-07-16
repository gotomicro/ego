package ejob

// Option 选项
type Option func(c *Container)

// WithName 设置Job的名称
func WithName(name string) Option {
	return func(c *Container) {
		c.config.Name = name
	}
}

// WithStartFunc 设置Job的函数
func WithStartFunc(startFunc func(ctx Context) error) Option {
	return func(c *Container) {
		c.config.startFunc = startFunc
	}
}

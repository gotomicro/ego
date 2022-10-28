package egovernor

// Option overrides a Container's default configuration.
type Option func(c *Container)

// WithHost 设置监控IP
func WithHost(host string) Option {
	return func(c *Container) {
		c.config.Host = host
	}
}

// WithPort 设置监控端口
func WithPort(port int) Option {
	return func(c *Container) {
		c.config.Port = port
	}
}

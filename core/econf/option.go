package econf

// Option 选项
type Option func(o *Container)

// WithTagName 设置Tag
func WithTagName(tag string) Option {
	return func(o *Container) {
		o.TagName = tag
	}
}

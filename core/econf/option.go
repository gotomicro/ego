package econf

// Option 选项
type Option func(o *Container)

// WithTagName 设置Tag
func WithTagName(tag ConfigType) Option {
	return func(o *Container) {
		o.TagName = string(tag)
	}
}

// WithWeaklyTypedInput 设置WeaklyTypedInput
func WithWeaklyTypedInput(weaklyTypedInput bool) Option {
	return func(o *Container) {
		o.WeaklyTypedInput = weaklyTypedInput
	}
}

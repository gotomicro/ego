package econf

// Option is an optional argument to Container.
type Option func(o *Container)

// WithTagName sets tagName when unmarshal raw config to struct.
func WithTagName(tag ConfigType) Option {
	return func(o *Container) {
		o.TagName = string(tag)
	}
}

// WithWeaklyTypedInput sets if allow weaklyTypedInput.
func WithWeaklyTypedInput(weaklyTypedInput bool) Option {
	return func(o *Container) {
		o.WeaklyTypedInput = weaklyTypedInput
	}
}

// WithSquash sets if allow squash tags in embedded struct
func WithSquash(squash bool) Option {
	return func(o *Container) {
		o.Squash = squash
	}
}

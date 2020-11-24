package conf

// GetOption ...
type (
	GetOption  func(o *GetOptions)
	GetOptions struct {
		TagName string
	}
)

var defaultGetOptions = GetOptions{
	TagName: "mapstructure",
}

// 设置Tag
func TagName(tag string) GetOption {
	return func(o *GetOptions) {
		o.TagName = tag
	}
}

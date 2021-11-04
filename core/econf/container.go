package econf

// Container 容器
type Container struct {
	TagName          string
	WeaklyTypedInput bool
}

var defaultContainer = Container{
	TagName:          "mapstructure",
	WeaklyTypedInput: false,
}

// GetOptionTagName 获取当前optionTag
func GetOptionTagName() string {
	return defaultContainer.TagName
}

// GetOptionWeaklyTypedInput 获取当前WeaklyTypedInput
func GetOptionWeaklyTypedInput() bool {
	return defaultContainer.WeaklyTypedInput
}

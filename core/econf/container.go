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

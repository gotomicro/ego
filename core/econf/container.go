package econf

// Container 容器
type Container struct {
	TagName string
}

var defaultContainer = Container{
	TagName: "mapstructure",
}

package econf

// Container defines a component instance.
type Container struct {
	TagName          string
	WeaklyTypedInput bool
	Squash           bool
}

var defaultContainer = Container{
	TagName:          "mapstructure",
	WeaklyTypedInput: false,
	Squash:           false,
}

// GetOptionTagName returns optionTag config of default container
func GetOptionTagName() string {
	return defaultContainer.TagName
}

// GetOptionWeaklyTypedInput returns WeaklyTypedInput config of default container
func GetOptionWeaklyTypedInput() bool {
	return defaultContainer.WeaklyTypedInput
}

// GetOptionSquash returns Squash config of default container
func GetOptionSquash() bool {
	return defaultContainer.Squash
}

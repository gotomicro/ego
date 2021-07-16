package ejob

// Config ...
type Config struct {
	Name string
	// context.Context 替换为 ejob.Context
	startFunc func(ctx Context) error
}

// defaultConfig 默认配置
func defaultConfig() *Config {
	return &Config{
		Name:      "",
		startFunc: nil,
	}
}

package ejob

// Config ...
type Config struct {
	Name      string
	startFunc func() error
	stopFunc  func() error
}

func DefaultConfig() *Config {
	return &Config{
		Name:      "",
		startFunc: nil,
		stopFunc:  nil,
	}
}

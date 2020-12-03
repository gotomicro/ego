package econf

import (
	"io"

	"github.com/davecgh/go-spew/spew"
)

// DataSource ...
type DataSource interface {
	Parse(addr string, watch bool)
	ReadConfig() ([]byte, error)
	IsConfigChanged() <-chan struct{}
	io.Closer
}

// Unmarshaller ...
type Unmarshaller = func([]byte, interface{}) error

var defaultConfiguration = New()

// OnChange 注册change回调函数
func OnChange(fn func(*Configuration)) {
	defaultConfiguration.OnChange(fn)
}

// LoadFromDataSource load configuration from data source
// if data source supports dynamic config, a monitor goroutinue
// would be
func LoadFromDataSource(ds DataSource, unmarshaller Unmarshaller, opts ...GetOption) error {
	return defaultConfiguration.LoadFromDataSource(ds, unmarshaller, opts...)
}

// Load loads configuration from provided provider with default defaultConfiguration.
func LoadFromReader(r io.Reader, unmarshaller Unmarshaller) error {
	return defaultConfiguration.LoadFromReader(r, unmarshaller)
}

// Apply ...
func Apply(conf map[string]interface{}) error {
	return defaultConfiguration.apply(conf)
}

// Reset resets all to default settings.
func Reset() {
	defaultConfiguration = New()
}

// Traverse ...
func Traverse(sep string) map[string]interface{} {
	return defaultConfiguration.traverse(sep)
}

// Debug ...
func Debug(sep string) {
	spew.Dump("Debug", Traverse(sep))
}

// Get returns an interface. For a specific value use one of the Get____ methods.
func Get(key string) interface{} {
	return defaultConfiguration.Get(key)
}

// Set set config value for key
func Set(key string, val interface{}) {
	defaultConfiguration.Set(key, val)
}

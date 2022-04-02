package econf

import (
	"io"

	"github.com/davecgh/go-spew/spew"
)

var (
	// ConfigTypeJSON ...
	ConfigTypeJSON ConfigType = "json"
	// ConfigTypeToml ...
	ConfigTypeToml ConfigType = "toml"
	// ConfigTypeYaml ...
	ConfigTypeYaml ConfigType = "yaml"
)

// ConfigType 配置类型
type ConfigType string

// DataSource ...
type DataSource interface {
	Parse(addr string, watch bool) ConfigType
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

// Sub return sub-configuration of defaultConfiguration
func Sub(key string) *Configuration {
	return defaultConfiguration.Sub(key)
}

// LoadFromDataSource load configuration from data source
// if data source supports dynamic config, a monitor goroutinue
// would be
func LoadFromDataSource(ds DataSource, unmarshaller Unmarshaller, opts ...Option) error {
	return defaultConfiguration.LoadFromDataSource(ds, unmarshaller, opts...)
}

// LoadFromReader loads configuration from provided provider with default defaultConfiguration.
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

// RawConfig 原始配置
func RawConfig() []byte {
	return defaultConfiguration.raw()
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
	_ = defaultConfiguration.Set(key, val)
}

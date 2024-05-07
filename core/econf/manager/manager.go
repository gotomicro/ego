package manager

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"

	"github.com/gotomicro/ego/core/econf"
)

var defaultScheme = "file"

var (
	// ErrInvalidDataSource defines an error that the scheme has been registered.
	ErrInvalidDataSource = errors.New("invalid data source, please make sure the scheme has been registered")
	// ErrInvalidUnmarshaller defines an error that unmarshaller is not exists.
	ErrInvalidUnmarshaller = errors.New("invalid unmarshaller, please make sure the config type is right")
	// ErrDefaultConfigNotExist defines an error than config not exists.
	ErrDefaultConfigNotExist = errors.New("default config not exist")
	registry                 map[string]econf.DataSource

	unmarshallers = map[econf.ConfigType]econf.Unmarshaller{
		econf.ConfigTypeJSON: json.Unmarshal,
		econf.ConfigTypeToml: toml.Unmarshal,
		econf.ConfigTypeYaml: yaml.Unmarshal,
	}
)

// DataSourceCreatorFunc represents a dataSource creator function.
type DataSourceCreatorFunc func() econf.DataSource

func init() {
	registry = make(map[string]econf.DataSource)
}

// Register registers a dataSource creator function to the registry.
func Register(scheme string, creator econf.DataSource) {
	registry[scheme] = creator
}

// NewDataSource constructs a new configuration provider by supplied config address.
func NewDataSource(configAddr string, watch bool) (econf.DataSource, econf.Unmarshaller, econf.ConfigType, error) {
	var scheme = defaultScheme
	urlObj, err := url.Parse(configAddr)
	if err == nil && len(urlObj.Scheme) > 1 {
		scheme = urlObj.Scheme
	}

	// 如果是默认file协议，判断下文件是否存在
	if scheme == defaultScheme {
		_, err := os.Stat(configAddr)
		if err != nil {
			return nil, nil, "", ErrDefaultConfigNotExist
		}
	}

	creatorFunc, exist := registry[scheme]
	if !exist {
		return nil, nil, "", ErrInvalidDataSource
	}
	tag := creatorFunc.Parse(configAddr, watch)

	parser, flag := unmarshallers[tag]
	if !flag {
		return nil, nil, "", ErrInvalidUnmarshaller
	}
	return creatorFunc, parser, tag, nil
}

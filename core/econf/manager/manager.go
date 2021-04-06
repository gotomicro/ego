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

var (
	// ErrInvalidDataSource defines an error that the scheme has been registered
	ErrInvalidDataSource = errors.New("invalid data source, please make sure the scheme has been registered")
	// ErrInvalidUnmarshaller ...
	ErrInvalidUnmarshaller = errors.New("invalid unmarshaller, please make sure the config type is right")
	// ErrDefaultConfigNotExist 默认配置不存在
	ErrDefaultConfigNotExist = errors.New("default config not exist")
	registry                 map[string]econf.DataSource
	// DefaultScheme 默认协议
	DefaultScheme = "file"

	//
	unmarshallerMap = map[econf.ConfigType]econf.Unmarshaller{
		econf.ConfigTypeJSON: json.Unmarshal,
		econf.ConfigTypeToml: toml.Unmarshal,
		econf.ConfigTypeYaml: yaml.Unmarshal,
	}
)

// DataSourceCreatorFunc represents a dataSource creator function
type DataSourceCreatorFunc func() econf.DataSource

func init() {
	registry = make(map[string]econf.DataSource)
}

// Register registers a dataSource creator function to the registry
func Register(scheme string, creator econf.DataSource) {
	registry[scheme] = creator
}

// NewDataSource 根据配置地址，创建数据源
func NewDataSource(configAddr string, watch bool) (econf.DataSource, econf.Unmarshaller, econf.ConfigType, error) {
	scheme := DefaultScheme
	urlObj, err := url.Parse(configAddr)
	if err == nil && len(urlObj.Scheme) > 1 {
		scheme = urlObj.Scheme
	}

	// 如果是默认file协议，判断下文件是否存在
	if scheme == DefaultScheme {
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

	parser, flag := unmarshallerMap[tag]
	if !flag {
		return nil, nil, "", ErrInvalidUnmarshaller
	}
	return creatorFunc, parser, tag, nil
}

package manager

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
)

var (
	// ErrInvalidDataSource defines an error that the scheme has been registered
	ErrInvalidDataSource = errors.New("invalid data source, please make sure the scheme has been registered")
	// ErrDefaultConfigNotExist 默认配置不存在
	ErrDefaultConfigNotExist = errors.New("default config not exist")
	registry                 map[string]econf.DataSource
	// DefaultScheme 默认协议
	DefaultScheme = "file"
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
func NewDataSource(configAddr string, watch bool) (econf.DataSource, econf.Unmarshaller, string, error) {
	// 如果配置为空，那么赋值默认配置
	if configAddr == "" {
		configAddr = eapp.EgoConfigPath()
	}

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

	creatorFunc.Parse(configAddr, watch)

	parser, tag := extParser(configAddr)

	return creatorFunc, parser, tag, nil
}

func extParser(configAddr string) (econf.Unmarshaller, string) {
	ext := filepath.Ext(configAddr)
	switch ext {
	case ".json":
		return json.Unmarshal, "json"
	case ".toml":
		return toml.Unmarshal, "toml"
	case ".yaml":
		return yaml.Unmarshal, "yaml"
	default:
		// TODO 处理configAddr为ETCD的情况？
		elog.EgoLogger.Panic("data source: invalid configuration type")
	}
	return nil, ""
}

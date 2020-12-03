package manager

import (
	"errors"
	"github.com/gotomicro/ego/core/eapp"
	"os"

	"github.com/gotomicro/ego/core/econf"
)

var (
	// ErrInvalidDataSource defines an error that the scheme has been registered
	ErrInvalidDataSource     = errors.New("invalid data source, please make sure the scheme has been registered")
	ErrDefaultConfigNotExist = errors.New("default config not exit")
	registry                 map[string]econf.DataSource
	DefaultScheme            = "file"
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

//NewDataSource ..
func NewDataSource(scheme, configAddr string, watch bool) (econf.DataSource, error) {
	if scheme == DefaultScheme && configAddr == eapp.EgoConfigPath() {
		_, err := os.Stat(configAddr)
		if err != nil {
			return nil, ErrDefaultConfigNotExist
		}
	}

	creatorFunc, exist := registry[scheme]
	if !exist {
		return nil, ErrInvalidDataSource
	}

	creatorFunc.Parse(configAddr, watch)
	return creatorFunc, nil
}

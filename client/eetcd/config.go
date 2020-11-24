package eetcd

import (
	"github.com/gotomicro/ego/core/util/xtime"
	"time"
)

// Config ...
type Config struct {
	Endpoints        []string      `json:"endpoints"`
	CertFile         string        `json:"certFile"`
	KeyFile          string        `json:"keyFile"`
	CaCert           string        `json:"caCert"`
	BasicAuth        bool          `json:"basicAuth"`
	UserName         string        `json:"userName"`
	Password         string        `json:"-"`
	ConnectTimeout   time.Duration `json:"connectTimeout"` // 连接超时时间
	Secure           bool          `json:"secure"`
	AutoSyncInterval time.Duration `json:"autoAsyncInterval"` // 自动同步member list的间隔
	TTL              int           // 单位：s
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		BasicAuth:      false,
		ConnectTimeout: xtime.Duration("5s"),
		Secure:         false,
	}
}

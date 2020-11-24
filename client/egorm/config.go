package egorm

import (
	"github.com/gotomicro/ego/core/util/xtime"
	"time"
)

// config options
type Config struct {
	Name            string        // 数据库名称
	Dialect         string        // 选择数据库种类，默认mysql
	DSN             string        // DSN地址: mysql://root:secret@tcp(127.0.0.1:3306)/mysql?timeout=20s&readTimeout=20s
	Debug           bool          // Debug开关，默认关闭
	MaxIdleConns    int           // 最大空闲连接数，默认10
	MaxOpenConns    int           // 最大活动连接数，默认100
	ConnMaxLifetime time.Duration // 连接的最大存活时间，默认300s
	OnDialError     string        // 创建连接的错误级别，=panic时，如果创建失败，立即panic，默认连接不上panic
	SlowThreshold   time.Duration // 慢日志阈值，默认500ms
	DialTimeout     time.Duration // 拨超时时间，默认1s
	DisableMetric   bool          // 关闭指标采集，默人false，也就是说默认采集监控
	DisableTrace    bool          // 关闭链路追踪，默认false，也就是说默认采集trace
	DetailSQL       bool          // 记录错误sql时,是否打印包含参数的完整sql语句，select * from aid = ?;
	interceptors    []Interceptor
	dsnCfg          *DSN
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		DSN:             "",
		Dialect:         "mysql",
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: xtime.Duration("300s"),
		OnDialError:     "panic",
		SlowThreshold:   xtime.Duration("500ms"),
		DialTimeout:     xtime.Duration("1s"),
		DisableMetric:   false,
		DisableTrace:    false,
	}
}

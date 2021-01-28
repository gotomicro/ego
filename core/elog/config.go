package elog

import (
	"fmt"
	"time"

	"github.com/gotomicro/ego/core/eapp"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config ...
type Config struct {
	Debug               bool          // 是否双写至文件控制日志输出到终端
	Level               string        // 日志初始等级，默认info级别
	Dir                 string        // [fileWriter]日志输出目录，默认logs
	Name                string        // [fileWriter]日志文件名称，默认框架日志mocro.sys，业务日志default.log
	MaxSize             int           // [fileWriter]日志输出文件最大长度，超过改值则截断，默认500M
	MaxAge              int           // [fileWriter]日志存储最大时间，默认最大保存天数为7天
	MaxBackup           int           // [fileWriter]日志存储最大数量，默认最大保存文件个数为10个
	RotateInterval      time.Duration // [fileWriter]日志轮转时间，默认1天
	EnableAddCaller     bool          // 是否添加调用者信息，默认不加调用者信息
	EnableAsync         bool          // 是否异步，默认异步
	FlushBufferSize     int           // 缓冲大小，默认256 * 1024B
	FlushBufferInterval time.Duration // 缓冲时间，默认5秒
	Writer              string        // 使用哪种Writer，默认使用fileWriter
	AliAccessKeyID      string        // [aliWriter]阿里云sls AKID
	AliAccessKeySecret  string        // [aliWriter]阿里云sls AKSecret
	AliEndpoint         string        // [aliWriter]阿里云sls endpoint
	AliProject          string        // [aliWriter]阿里云sls Project名称
	AliLogstore         string        // [aliWriter]阿里云sls logstore名称
	AliTopic            string        // [aliWriter]阿里云sls logstore名称

	fields        []zap.Field // 日志初始化字段
	callerSkip    int
	core          zapcore.Core
	encoderConfig *zapcore.EncoderConfig
	configKey     string
}

const (
	writerRotateFile = "file"
	writerAliSLS     = "ali"
)

// filename ...
func (config *Config) filename() string {
	return fmt.Sprintf("%s/%s", config.Dir, config.Name)
}

// DefaultConfig ...
func DefaultConfig() *Config {
	dir := "./logs"
	if eapp.EgoLogPath() != "" {
		dir = eapp.EgoLogPath()
	}
	return &Config{
		Name:                DefaultLoggerName,
		Dir:                 dir,
		Level:               "info",
		FlushBufferSize:     defaultBufferSize,
		FlushBufferInterval: defaultFlushInterval,
		MaxSize:             500, // 500M
		MaxAge:              7,   // 1 day
		MaxBackup:           10,  // 10 backup
		RotateInterval:      24 * time.Hour,
		callerSkip:          1,
		EnableAddCaller:     false,
		EnableAsync:         true,
		encoderConfig:       DefaultZapConfig(),
		Writer:              writerRotateFile,
	}
}

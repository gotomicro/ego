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
	Debug                     bool          // 是否双写至文件控制日志输出到终端
	Level                     string        // 日志初始等级，默认info级别
	Dir                       string        // [fileWriter]日志输出目录，默认logs
	Name                      string        // [fileWriter]日志文件名称，默认框架日志ego.sys，业务日志default.log
	MaxSize                   int           // [fileWriter]日志输出文件最大长度，超过改值则截断，默认500M
	MaxAge                    int           // [fileWriter]日志存储最大时间，默认最大保存天数为7天
	MaxBackup                 int           // [fileWriter]日志存储最大数量，默认最大保存文件个数为10个
	RotateInterval            time.Duration // [fileWriter]日志轮转时间，默认1天
	EnableAddCaller           bool          // 是否添加调用者信息，默认不加调用者信息
	EnableAsync               bool          // 是否异步，默认异步
	FlushBufferSize           int           // 缓冲大小，默认256 * 1024B
	FlushBufferInterval       time.Duration // 缓冲时间，默认5秒
	Writer                    string        // 使用哪种Writer，可选[file|ali|stderr]，默认file
	AliAccessKeyID            string        // [aliWriter]阿里云sls AKID，必填
	AliAccessKeySecret        string        // [aliWriter]阿里云sls AKSecret，必填
	AliEndpoint               string        // [aliWriter]阿里云sls endpoint，必填
	AliProject                string        // [aliWriter]阿里云sls Project名称，必填
	AliLogstore               string        // [aliWriter]阿里云sls logstore名称，必填
	AliAPIBulkSize            int           // [aliWriter]阿里云sls API单次请求发送最大日志条数，最少256条，默认256条
	AliAPITimeout             time.Duration // [aliWriter]阿里云sls API接口超时，默认3秒
	AliAPIRetryCount          int           // [aliWriter]阿里云sls API接口重试次数，默认3次
	AliAPIRetryWaitTime       time.Duration // [aliWriter]阿里云sls API接口重试默认等待间隔，默认1秒
	AliAPIRetryMaxWaitTime    time.Duration // [aliWriter]阿里云sls API接口重试最大等待间隔，默认3秒
	AliAPIMaxIdleConnsPerHost int           // [aliWriter]阿里云sls 单个Host HTTP最大空闲连接数，应当大于AliApiMaxIdleConns
	AliAPIMaxIdleConns        int           // [aliWriter]阿里云sls HTTP最大空闲连接数
	AliAPIIdleConnTimeout     time.Duration // [aliWriter]阿里云sls HTTP空闲连接保活时间

	fields        []zap.Field // 日志初始化字段
	CallerSkip    int
	encoderConfig *zapcore.EncoderConfig
}

const (
	writerRotateFile = "file"
	writerAliSLS     = "ali"
	writerStderr     = "stderr"
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
		Name:                      DefaultLoggerName,
		Dir:                       dir,
		Level:                     "info",
		FlushBufferSize:           defaultBufferSize,
		FlushBufferInterval:       defaultFlushInterval,
		MaxSize:                   500, // 500M
		MaxAge:                    7,   // 1 day
		MaxBackup:                 10,  // 10 backup
		RotateInterval:            24 * time.Hour,
		CallerSkip:                1,
		EnableAddCaller:           false,
		EnableAsync:               true,
		encoderConfig:             defaultZapConfig(),
		Writer:                    writerRotateFile,
		AliAPIBulkSize:            256,
		AliAPITimeout:             3 * time.Second,
		AliAPIRetryCount:          3,
		AliAPIRetryWaitTime:       1 * time.Second,
		AliAPIRetryMaxWaitTime:    3 * time.Second,
		AliAPIMaxIdleConnsPerHost: 20,
		AliAPIMaxIdleConns:        25,
		AliAPIIdleConnTimeout:     30 * time.Second,
	}
}

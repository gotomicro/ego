package elog

import (
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/logrotate"
)

type rotateWriterBuilder struct{}

type rotateWriter struct {
	zapcore.Core
	io.Closer
}

var _ WriterBuilder = &rotateWriterBuilder{}

// config ...
type config struct {
	MaxSize             int           // [fileWriter]日志输出文件最大长度，超过改值则截断，默认500M
	MaxAge              int           // [fileWriter]日志存储最大时间，默认最大保存天数为7天
	MaxBackup           int           // [fileWriter]日志存储最大数量，默认最大保存文件个数为10个
	RotateInterval      time.Duration // [fileWriter]日志轮转时间，默认0，不做时间轮转，影响性能
	FlushBufferSize     int           // 缓冲大小，默认256 * 1024B
	FlushBufferInterval time.Duration // 缓冲时间，默认5秒
}

func defaultConfig() *config {
	return &config{
		MaxSize:             500, // 500M
		MaxAge:              7,   // 7 day
		MaxBackup:           10,  // 10 backup
		RotateInterval:      0,
		FlushBufferSize:     256 * 1024,
		FlushBufferInterval: 5 * time.Second,
	}
}

const (
	writerRotateLogger = "file"
)

func (*rotateWriterBuilder) Scheme() string {
	return writerRotateLogger
}

// Build constructs a zapcore.Core with stderr syncer
func (r *rotateWriterBuilder) Build(key string, commonConfig *Config) Writer {
	c := defaultConfig()
	if err := econf.UnmarshalKey(key, &c); err != nil {
		panic(err)
	}
	// NewRotateFileCore constructs a zapcore.Core with rotate file syncer
	// Debug output to console and file by default
	cf := noopCloseFunc
	var ws = zapcore.AddSync(&logrotate.Logger{
		Filename:       commonConfig.Filename(),
		MaxSize:        c.MaxSize,
		MaxAge:         c.MaxAge,
		MaxBackups:     c.MaxBackup,
		LocalTime:      true,
		Compress:       false,
		RotateInterval: c.RotateInterval,
	})

	if commonConfig.Debug {
		ws = zap.CombineWriteSyncers(os.Stdout, ws)
	}
	if commonConfig.EnableAsync {
		ws, cf = bufferWriteSyncer(ws, c.FlushBufferSize, c.FlushBufferInterval)
	}
	w := &rotateWriter{}
	w.Closer = CloseFunc(cf)
	w.Core = zapcore.NewCore(
		func() zapcore.Encoder {
			if commonConfig.Debug {
				return zapcore.NewConsoleEncoder(*commonConfig.EncoderConfig())
			}
			return zapcore.NewJSONEncoder(*commonConfig.EncoderConfig())
		}(),
		ws,
		commonConfig.AtomicLevel(),
	)
	return w
}

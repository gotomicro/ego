package elog

import (
	"io"

	"github.com/gotomicro/ego/core/elog/rotate"
)

func newRotate(config *Config) io.Writer {
	rotateLog := rotate.NewLogger()
	rotateLog.Filename = config.Filename()
	rotateLog.MaxSize = config.MaxSize // MB
	rotateLog.MaxAge = config.MaxAge   // days
	rotateLog.MaxBackups = config.MaxBackup
	rotateLog.Interval = config.RotateInterval
	rotateLog.LocalTime = true
	rotateLog.Compress = false
	return rotateLog
}

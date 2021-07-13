package elog

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// fileDataSource file provider.
type stderrLogger struct{}

// Load constructs a zapcore.Core with stderr syncer
func (*stderrLogger) Load(key string, commonConfig *Config, lv zap.AtomicLevel) (zapcore.Core, CloseFunc) {
	// Debug output to console and file by default
	return zapcore.NewCore(zapcore.NewJSONEncoder(*commonConfig.EncoderConfig()), os.Stderr, lv), noopCloseFunc
}

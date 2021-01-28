package ali

import (
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewCore(t *testing.T) {
	lv := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	core, cancel := NewCore(
		WithLevelEnabler(lv),
		WithFlushBufferInterval(5*time.Second),
		WithFlushBufferInterval(256*1024),
	)
	defer cancel()
	zapLogger := zap.New(core).With(zap.String("prefix", "PREFIX"))
	zapLogger.Error("my message", zap.String("aaa", "AAA"), zap.Int("bbb", 111), zap.String("aaa", "CCC"))
}

package elog

import (
	"io"
	"os"

	"go.uber.org/zap/zapcore"
)

const (
	writerStderr = "stderr"
)

var _ WriterBuilder = &stderrWriterBuilder{}

// fileDataSource file Provider.
type stderrWriterBuilder struct{}

type stderrWriter struct {
	zapcore.Core
	io.Closer
}

func (s *stderrWriterBuilder) Build(key string, c *Config) Writer {
	// Debug output to console and file by default
	w := &stderrWriter{}
	w.Core = zapcore.NewCore(zapcore.NewJSONEncoder(*c.EncoderConfig()), os.Stderr, c.AtomicLevel())
	w.Closer = CloseFunc(noopCloseFunc)
	return w
}

func (*stderrWriterBuilder) Scheme() string {
	return writerStderr
}

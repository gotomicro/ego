package elog

import (
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

const (
	writerStdout = "stdout"
)

var _ WriterBuilder = &stdoutWriterBuilder{}

// fileDataSource file Provider.
type stdoutWriterBuilder struct{}

type stdoutWriter struct {
	zapcore.Core
	io.Closer
}

func (s *stdoutWriterBuilder) Build(key string, c *Config) Writer {
	// Debug output to console and file by default
	w := &stdoutWriter{}
	w.Core = zapcore.NewCore(zapcore.NewJSONEncoder(*c.EncoderConfig()), os.Stdout, c.AtomicLevel())
	w.Closer = CloseFunc(noopCloseFunc)
	return w
}

func (*stdoutWriterBuilder) Scheme() string {
	return writerStdout
}

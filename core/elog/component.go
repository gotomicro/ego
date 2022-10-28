package elog

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gotomicro/ego/core/econf"
)

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zap.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zap.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zap.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-Level logs.
	ErrorLevel = zap.ErrorLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = zap.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zap.FatalLevel
)

type (
	// Field ...
	Field = zap.Field
	// Level ...
	Level = zapcore.Level
	// Component defines ego logger component wraps zap.Logger and zap.SugaredLogger
	Component struct {
		name          string
		desugar       *zap.Logger
		lv            *zap.AtomicLevel
		config        *Config
		sugar         *zap.SugaredLogger
		asyncStopFunc func() error
	}
)

var (
	// String alias for zap.String
	String = zap.String
	// Any alias for zap.Any
	Any = zap.Any
	// Int64 alias for zap.Int64
	Int64 = zap.Int64
	// Int alias for zap.Int
	Int = zap.Int
	// Int32 alias for zap.Int32
	Int32 = zap.Int32
	// Uint alias for zap.Uint
	Uint = zap.Uint
	// Duration alias for zap.Duration
	Duration = zap.Duration
	// Durationp alias for zap.Duration
	Durationp = zap.Durationp
	// Object alias for zap.Object
	Object = zap.Object
	// Namespace alias for zap.Namespace
	Namespace = zap.Namespace
	// Reflect alias for zap.Reflect
	Reflect = zap.Reflect
	// Skip alias for zap.Skip()
	Skip = zap.Skip()
	// ByteString alias for zap.ByteString
	ByteString = zap.ByteString
)

func newLogger(name string, key string, config *Config) *Component {
	zapOptions := make([]zap.Option, 0)
	zapOptions = append(zapOptions, zap.AddStacktrace(zap.DPanicLevel))
	if config.EnableAddCaller {
		zapOptions = append(zapOptions, zap.AddCaller(), zap.AddCallerSkip(config.CallerSkip))
	}
	if len(config.fields) > 0 {
		zapOptions = append(zapOptions, zap.Fields(config.fields...))
	}

	// unmarshals the text to an AtomicLevel.
	if err := config.al.UnmarshalText([]byte(config.Level)); err != nil {
		panic(err)
	}

	// sets core to default zap.Core if not configured.
	if config.core == nil {
		w := Provider(config.Writer).Build(key, config)
		config.core = w
		config.asyncStopFunc = w.Close
	}

	zapLogger := zap.New(config.core, zapOptions...)
	l := &Component{
		desugar:       zapLogger,
		lv:            &config.al,
		config:        config,
		sugar:         zapLogger.Sugar(),
		name:          name,
		asyncStopFunc: config.asyncStopFunc,
	}

	// 如果名字不为空，加载动态配置
	if l.name != "" {
		l.AutoLevel(name + ".level")
	}
	return l

}

// ZapLogger returns *zap.Logger
func (logger *Component) ZapLogger() *zap.Logger {
	return logger.desugar
}

// ZapSugaredLogger returns *zap.SugaredLogger
func (logger *Component) ZapSugaredLogger() *zap.SugaredLogger {
	return logger.sugar
}

// AutoLevel ...
func (logger *Component) AutoLevel(confKey string) {
	econf.OnChange(func(config *econf.Configuration) {
		lvText := strings.ToLower(config.GetString(confKey))
		if lvText != "" {
			logger.Info("update level", String("level", lvText), String("name", logger.config.Name))
			_ = logger.lv.UnmarshalText([]byte(lvText))
		}
	})
}

// SetLevel ...
func (logger *Component) SetLevel(lv Level) {
	logger.lv.SetLevel(lv)
}

// Flush ...
// When use os.Stdout or os.Stderr as zapcore.WriteSyncer
// logger.desugar.Sync() maybe return an error like this: 'sync /dev/stdout: The handle is invalid.'
// Because os.Stdout and os.Stderr is a non-normal file, maybe not support 'fsync' in different os platform
// So ignored Sync() return value
// About issues: https://github.com/uber-go/zap/issues/328
// About 'fsync': https://man7.org/linux/man-pages/man2/fsync.2.html
func (logger *Component) Flush() error {
	if logger.asyncStopFunc != nil {
		if err := logger.asyncStopFunc(); err != nil {
			return err
		}
	}

	_ = logger.desugar.Sync()
	return nil
}

// IsDebugMode ...
func (logger *Component) IsDebugMode() bool {
	return logger.config.Debug
}

// Debug ...
func (logger *Component) Debug(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.desugar.Debug(msg, fields...)
}

// Debugw ...
// Deprecated: Will be removed in future versions, use *Component.Debug instead.
func (logger *Component) Debugw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Debugw(msg, keysAndValues...)
}

// Debugf ...
// Deprecated: Will be removed in future versions, use *Component.Debug instead.
func (logger *Component) Debugf(template string, args ...interface{}) {
	logger.sugar.Debugf(template, args...)
}

// Info ...
func (logger *Component) Info(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.desugar.Info(msg, fields...)
}

// Infow ...
// Deprecated: Will be removed in future versions, use *Component.Info instead.
func (logger *Component) Infow(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Infow(msg, keysAndValues...)
}

// Infof ...
// Deprecated: Will be removed in future versions, use *Component.Info instead.
func (logger *Component) Infof(template string, args ...interface{}) {
	logger.sugar.Infof(template, args...)
}

// Warn ...
func (logger *Component) Warn(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.desugar.Warn(msg, fields...)
}

// Warnw ...
// Deprecated: Will be removed in future versions, use *Component.Warn instead.
func (logger *Component) Warnw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Warnw(msg, keysAndValues...)
}

// Warnf ...
// Deprecated: Will be removed in future versions, use *Component.Warn instead.
func (logger *Component) Warnf(template string, args ...interface{}) {
	logger.sugar.Warnf(template, args...)
}

// Error ...
func (logger *Component) Error(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.desugar.Error(msg, fields...)
}

// Errorw ...
// Deprecated: Will be removed in future versions, use *Component.Error instead.
func (logger *Component) Errorw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Errorw(msg, keysAndValues...)
}

// Errorf ...
// Deprecated: Will be removed in future versions, use *Component.Error instead.
func (logger *Component) Errorf(template string, args ...interface{}) {
	logger.sugar.Errorf(template, args...)
}

// Panic ...
func (logger *Component) Panic(msg string, fields ...Field) {
	panicDetail(msg, fields...)
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.desugar.Panic(msg, fields...)
}

// Panicw ...
// Deprecated: Will be removed in future versions, use *Component.Panic instead.
func (logger *Component) Panicw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Panicw(msg, keysAndValues...)
}

// Panicf ...
// Deprecated: Will be removed in future versions, use *Component.Panic instead.
func (logger *Component) Panicf(template string, args ...interface{}) {
	logger.sugar.Panicf(template, args...)
}

// DPanic ...
func (logger *Component) DPanic(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		panicDetail(msg, fields...)
		msg = normalizeMessage(msg)
	}
	logger.desugar.DPanic(msg, fields...)
}

// DPanicw ...
// Deprecated: Will be removed in future versions, use *Component.DPanic instead.
func (logger *Component) DPanicw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.DPanicw(msg, keysAndValues...)
}

// DPanicf ...
// Deprecated: Will be removed in future versions, use *Component.DPanic instead.
func (logger *Component) DPanicf(template string, args ...interface{}) {
	logger.sugar.DPanicf(template, args...)
}

// Fatal ...
func (logger *Component) Fatal(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		panicDetail(msg, fields...)
		// msg = normalizeMessage(msg)
		return
	}
	logger.desugar.Fatal(msg, fields...)
}

// Fatalw ...
// Deprecated: Will be removed in future versions, use *Component.Fatal instead.
func (logger *Component) Fatalw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Fatalw(msg, keysAndValues...)
}

// Fatalf ...
// Deprecated: Will be removed in future versions, use *Component.Fatal instead.
func (logger *Component) Fatalf(template string, args ...interface{}) {
	logger.sugar.Fatalf(template, args...)
}

// With creates a child logger
func (logger *Component) With(fields ...Field) *Component {
	desugarLogger := logger.desugar.With(fields...)
	return &Component{
		desugar: desugarLogger,
		lv:      logger.lv,
		sugar:   desugarLogger.Sugar(),
		config:  logger.config,
	}
}

// WithCallerSkip ...
func (logger *Component) WithCallerSkip(callerSkip int, fields ...Field) *Component {
	logger.config.CallerSkip = callerSkip
	desugarLogger := logger.desugar.WithOptions(zap.AddCallerSkip(callerSkip)).With(fields...)
	return &Component{
		desugar: desugarLogger,
		lv:      logger.lv,
		sugar:   desugarLogger.Sugar(),
		config:  logger.config,
	}
}

// ConfigDir returns log directory path if a fileWriter logger is set.
func (logger *Component) ConfigDir() string {
	return logger.config.Dir
}

// ConfigName returns logger name.
func (logger *Component) ConfigName() string {
	return logger.config.Name
}

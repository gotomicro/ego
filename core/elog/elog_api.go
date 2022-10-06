package elog

// PackageName 包名
const PackageName = "core.elog"

// DefaultLogger defines default logger, it's usually used in application business logic
var DefaultLogger *Component

// EgoLogger defines ego framework logger, it's for ego framework only
var EgoLogger *Component

func init() {
	registry = make(map[string]WriterBuilder)
	Register(&stderrWriterBuilder{})
	Register(&rotateWriterBuilder{})
	DefaultLogger = DefaultContainer().Build(WithFileName(DefaultLoggerName))
	EgoLogger = DefaultContainer().Build(WithFileName(EgoLoggerName))
}

// Info ...
func Info(msg string, fields ...Field) {
	DefaultLogger.Info(msg, fields...)
}

// Debug ...
func Debug(msg string, fields ...Field) {
	DefaultLogger.Debug(msg, fields...)
}

// Warn ...
func Warn(msg string, fields ...Field) {
	DefaultLogger.Warn(msg, fields...)
}

// Error ...
func Error(msg string, fields ...Field) {
	DefaultLogger.Error(msg, fields...)
}

// Panic ...
func Panic(msg string, fields ...Field) {
	DefaultLogger.Panic(msg, fields...)
}

// DPanic ...
func DPanic(msg string, fields ...Field) {
	DefaultLogger.DPanic(msg, fields...)
}

// Fatal ...
func Fatal(msg string, fields ...Field) {
	DefaultLogger.Fatal(msg, fields...)
}

// Debugw ...
// Deprecated: Will be removed in future versions, use Debug instead.
func Debugw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Debugw(msg, keysAndValues...)
}

// Infow ...
// Deprecated: Will be removed in future versions, use Info instead.
func Infow(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Infow(msg, keysAndValues...)
}

// Warnw ...
// Deprecated: Will be removed in future versions, use Warn instead.
func Warnw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Warnw(msg, keysAndValues...)
}

// Errorw ...
// Deprecated: Will be removed in future versions, use Error instead.
func Errorw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Errorw(msg, keysAndValues...)
}

// Panicw ...
// Deprecated: Will be removed in future versions, use Panic instead.
func Panicw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Panicw(msg, keysAndValues...)
}

// DPanicw ...
// Deprecated: Will be removed in future versions, use DPanic instead.
func DPanicw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.DPanicw(msg, keysAndValues...)
}

// Fatalw ...
// Deprecated: Will be removed in future versions, use Fatal instead.
func Fatalw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Fatalw(msg, keysAndValues...)
}

// Debugf ...
// Deprecated: Will be removed in future versions, use Debug instead.
func Debugf(msg string, args ...interface{}) {
	DefaultLogger.Debugf(msg, args...)
}

// Infof ...
// Deprecated: Will be removed in future versions, use Info instead.
func Infof(msg string, args ...interface{}) {
	DefaultLogger.Infof(msg, args...)
}

// Warnf ...
// Deprecated: Will be removed in future versions, use Warn instead.
func Warnf(msg string, args ...interface{}) {
	DefaultLogger.Warnf(msg, args...)
}

// Errorf ...
// Deprecated: Will be removed in future versions, use Error instead.
func Errorf(msg string, args ...interface{}) {
	DefaultLogger.Errorf(msg, args...)
}

// Panicf ...
// Deprecated: Will be removed in future versions, use Panic instead.
func Panicf(msg string, args ...interface{}) {
	DefaultLogger.Panicf(msg, args...)
}

// DPanicf ...
// Deprecated: Will be removed in future versions, use DPanic instead.
func DPanicf(msg string, args ...interface{}) {
	DefaultLogger.DPanicf(msg, args...)
}

// Fatalf ...
// Deprecated: Will be removed in future versions, use Fatal instead.
func Fatalf(msg string, args ...interface{}) {
	DefaultLogger.Fatalf(msg, args...)
}

// With ...
func With(fields ...Field) *Component {
	return DefaultLogger.With(fields...)
}

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

// Deprecated: Will be removed in future versions, use Debug instead.
// Debugw ...
func Debugw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Debugw(msg, keysAndValues...)
}

// Deprecated: Will be removed in future versions, use Info instead.
// Infow ...
func Infow(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Infow(msg, keysAndValues...)
}

// Deprecated: Will be removed in future versions, use Warn instead.
// Warnw ...
func Warnw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Warnw(msg, keysAndValues...)
}

// Deprecated: Will be removed in future versions, use Error instead.
// Errorw ...
func Errorw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Errorw(msg, keysAndValues...)
}

// Deprecated: Will be removed in future versions, use Panic instead.
// Panicw ...
func Panicw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Panicw(msg, keysAndValues...)
}

// Deprecated: Will be removed in future versions, use DPanic instead.
// DPanicw ...
func DPanicw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.DPanicw(msg, keysAndValues...)
}

// Deprecated: Will be removed in future versions, use Fatal instead.
// Fatalw ...
func Fatalw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Fatalw(msg, keysAndValues...)
}

// Deprecated: Will be removed in future versions, use Debug instead.
// Debugf ...
func Debugf(msg string, args ...interface{}) {
	DefaultLogger.Debugf(msg, args...)
}

// Deprecated: Will be removed in future versions, use Info instead.
// Infof ...
func Infof(msg string, args ...interface{}) {
	DefaultLogger.Infof(msg, args...)
}

// Deprecated: Will be removed in future versions, use Warn instead.
// Warnf ...
func Warnf(msg string, args ...interface{}) {
	DefaultLogger.Warnf(msg, args...)
}

// Deprecated: Will be removed in future versions, use Error instead.
// Errorf ...
func Errorf(msg string, args ...interface{}) {
	DefaultLogger.Errorf(msg, args...)
}

// Deprecated: Will be removed in future versions, use Panic instead.
// Panicf ...
func Panicf(msg string, args ...interface{}) {
	DefaultLogger.Panicf(msg, args...)
}

// Deprecated: Will be removed in future versions, use DPanic instead.
// DPanicf ...
func DPanicf(msg string, args ...interface{}) {
	DefaultLogger.DPanicf(msg, args...)
}

// Deprecated: Will be removed in future versions, use Fatal instead.
// Fatalf ...
func Fatalf(msg string, args ...interface{}) {
	DefaultLogger.Fatalf(msg, args...)
}

// With ...
func With(fields ...Field) *Component {
	return DefaultLogger.With(fields...)
}

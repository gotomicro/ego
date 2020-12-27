package elog

const PackageName = "core.elog"

// DefaultLogger default logger
// Biz Log
// debug=true as default, will be
var DefaultLogger *Component

// frame logger
var EgoLogger *Component

func init() {
	DefaultLogger = DefaultContainer().Build()
	EgoLogger = DefaultContainer().Build(WithFileName(EgoLoggerName))
}

const (
	DefaultLoggerName = "default.log" // 默认业务日志名称
	EgoLoggerName     = "ego.sys"     // 框架日志名称
)

// Auto ...
func Auto(err error) Func {
	if err != nil {
		return DefaultLogger.With(FieldErr(err)).Error
	}
	return DefaultLogger.Info
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
func Debugw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Debugw(msg, keysAndValues...)
}

// Infow ...
func Infow(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Infow(msg, keysAndValues...)
}

// Warnw ...
func Warnw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Warnw(msg, keysAndValues...)
}

// Errorw ...
func Errorw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Errorw(msg, keysAndValues...)
}

// Panicw ...
func Panicw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Panicw(msg, keysAndValues...)
}

// DPanicw ...
func DPanicw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.DPanicw(msg, keysAndValues...)
}

// Fatalw ...
func Fatalw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Fatalw(msg, keysAndValues...)
}

// Debugf ...
func Debugf(msg string, args ...interface{}) {
	DefaultLogger.Debugf(msg, args...)
}

// Infof ...
func Infof(msg string, args ...interface{}) {
	DefaultLogger.Infof(msg, args...)
}

// Warnf ...
func Warnf(msg string, args ...interface{}) {
	DefaultLogger.Warnf(msg, args...)
}

// Errorf ...
func Errorf(msg string, args ...interface{}) {
	DefaultLogger.Errorf(msg, args...)
}

// Panicf ...
func Panicf(msg string, args ...interface{}) {
	DefaultLogger.Panicf(msg, args...)
}

// DPanicf ...
func DPanicf(msg string, args ...interface{}) {
	DefaultLogger.DPanicf(msg, args...)
}

// Fatalf ...
func Fatalf(msg string, args ...interface{}) {
	DefaultLogger.Fatalf(msg, args...)
}

// Log ...
func (fn Func) Log(msg string, fields ...Field) {
	fn(msg, fields...)
}

// With ...
func With(fields ...Field) *Component {
	return DefaultLogger.With(fields...)
}

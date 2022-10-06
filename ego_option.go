package ego

import (
	"os"
	"time"
)

// Option 可选项
type Option func(a *Ego)

// WithHang 是否允许系统悬挂起来，0 表示不悬挂， 1 表示悬挂。目的是一些脚本操作的时候，不想主线程停止
func WithHang(flag bool) Option {
	return func(a *Ego) {
		a.opts.hang = flag
	}
}

// WithDisableBanner 禁止banner
func WithDisableBanner(disableBanner bool) Option {
	return func(a *Ego) {
		a.opts.disableBanner = disableBanner
	}
}

// WithArguments 传入arguments
func WithArguments(arguments []string) Option {
	return func(a *Ego) {
		a.opts.arguments = arguments
	}
}

// WithDisableFlagConfig 禁止config
func WithDisableFlagConfig(disableFlagConfig bool) Option {
	return func(a *Ego) {
		a.opts.disableFlagConfig = disableFlagConfig
	}
}

// WithConfigPrefix 设置配置前缀
func WithConfigPrefix(configPrefix string) Option {
	return func(a *Ego) {
		a.opts.configPrefix = configPrefix
	}
}

// WithBeforeStopClean 设置运行前清理
func WithBeforeStopClean(fns ...func() error) Option {
	return func(a *Ego) {
		a.opts.beforeStopClean = append(a.opts.beforeStopClean, fns...)
	}
}

// WithAfterStopClean 设置运行后清理
func WithAfterStopClean(fns ...func() error) Option {
	return func(a *Ego) {
		a.opts.afterStopClean = append(a.opts.afterStopClean, fns...)
	}
}

// WithStopTimeout 设置停止的超时时间
func WithStopTimeout(timeout time.Duration) Option {
	return func(e *Ego) {
		e.opts.stopTimeout = timeout
	}
}

// WithShutdownSignal 设置停止信号量
func WithShutdownSignal(signals ...os.Signal) Option {
	return func(e *Ego) {
		e.opts.shutdownSignals = append(e.opts.shutdownSignals, signals...)
	}
}

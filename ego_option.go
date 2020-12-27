package ego

import (
	"os"
	"time"
)

type Option func(a *ego)

// 是否允许系统悬挂起来，0 表示不悬挂， 1 表示悬挂。目的是一些脚本操作的时候，不想主线程停止
func WithHang(flag bool) Option {
	return func(a *ego) {
		a.opts.hang = flag
	}
}

func WithDisableBanner(disableBanner bool) Option {
	return func(a *ego) {
		a.opts.disableBanner = disableBanner
	}
}

func WithConfigPrefix(configPrefix string) Option {
	return func(a *ego) {
		a.opts.configPrefix = configPrefix
	}
}

// 设置运行前清理
func WithBeforeStopClean(fns ...func() error) Option {
	return func(a *ego) {
		a.opts.beforeStopClean = append(a.opts.beforeStopClean, fns...)
	}
}

// 设置运行后清理
func WithAfterStopClean(fns ...func() error) Option {
	return func(a *ego) {
		a.opts.afterStopClean = append(a.opts.afterStopClean, fns...)
	}
}

func WithStopTimeout(timeout time.Duration) Option {
	return func(e *ego) {
		e.opts.stopTimeout = timeout
	}
}

func WithShutdownSignal(signals ...os.Signal) Option {
	return func(e *ego) {
		e.opts.shutdownSignals = append(e.opts.shutdownSignals, signals...)
	}
}

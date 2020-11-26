package ego

import "github.com/gotomicro/ego/core/registry"

type Option func(a *ego)

// SetRegistry set customize registry
func WithRegistry(reg registry.Registry) Option {
	return func(a *ego) {
		a.registerer = reg
	}
}

// 是否允许系统悬挂起来，0 表示不悬挂， 1 表示悬挂。目的是一些脚本操作的时候，不想主线程停止
func WithHang(flag bool) Option {
	return func(a *ego) {
		a.hang = flag
	}
}

// 设置运行前清理
func WithBeforeStopClean(fns ...func() error) Option {
	return func(a *ego) {
		a.beforeStopClean = fns
	}
}

// 设置运行后清理
func WithAfterStopClean(fns ...func() error) Option {
	return func(a *ego) {
		a.afterStopClean = fns
	}
}

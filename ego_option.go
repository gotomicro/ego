package ego

type Option func(a *ego)

// 是否允许系统悬挂起来，0 表示不悬挂， 1 表示悬挂。目的是一些脚本操作的时候，不想主线程停止
func WithHang(flag bool) Option {
	return func(a *ego) {
		a.hang = flag
	}
}

// 设置运行前清理
func WithBeforeStopClean(fns ...func() error) Option {
	return func(a *ego) {
		a.beforeStopClean = append(a.beforeStopClean, fns...)
	}
}

// 设置运行后清理
func WithAfterStopClean(fns ...func() error) Option {
	return func(a *ego) {
		a.afterStopClean = append(a.afterStopClean, fns...)
	}
}

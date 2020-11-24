package xgo

var (
//_logger = xlog.EgoLogger.With(zap.String("mod", "xgo"))
)

func try(fn func() error, cleaner func()) (ret error) {
	if cleaner != nil {
		defer cleaner()
	}
	return fn()
}

func tryIgnoreFnReturn(fn func(), cleaner func()) (ret error) {
	if cleaner != nil {
		defer cleaner()
	}
	fn()
	return nil
}

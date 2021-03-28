package xgo

// Go goroutine
func Go(fn func()) {
	go fn()
}

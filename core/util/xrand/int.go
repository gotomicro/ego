package xrand

import (
	"math/rand"
	"time"
)

// Int63n implements rand.Int63n on the grpcrand global source.
func Int63n(n int64) int64 {
	mu.Lock()
	res := r.Int63n(n)
	mu.Unlock()
	return res
}

// Intn implements rand.Intn on the grpcrand global source.
func Intn(n int) int {
	mu.Lock()
	res := r.Intn(n)
	mu.Unlock()
	return res
}

// Float64 implements rand.Float64 on the grpcrand global source.
func Float64() float64 {
	mu.Lock()
	res := r.Float64()
	mu.Unlock()
	return res
}

// Shuffle ...
func Shuffle(length int, fn func(i, j int)) {
	xr := rand.New(rand.NewSource(time.Now().UnixNano()))
	xr.Shuffle(length, fn)
}

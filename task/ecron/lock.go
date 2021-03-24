package ecron

import (
	"context"
	"sync"
	"time"
)

// Lock ...
// implementations:
//		Redis: [ecronlock](github.com/gotomicro/eredis@v0.2.0+)
//
type Lock interface {
	Lock(ctx context.Context, ttl time.Duration) error
	Unlock(ctx context.Context) error
	Refresh(ctx context.Context, ttl time.Duration) error
}

type mockLock struct {
	key    string
	mtx    *sync.Mutex
	locked bool
}

func (m *mockLock) WithKey(key string) *mockLock {
	return &mockLock{
		key: key,
		mtx: m.mtx,
	}
}

func (m *mockLock) Lock(ctx context.Context, ttl time.Duration) error {
	m.mtx.Lock()
	m.locked = true
	return nil
}

func (m *mockLock) Unlock(ctx context.Context) error {
	if m.locked {
		m.mtx.Unlock()
		m.locked = false
	}
	return nil
}

func (m *mockLock) Refresh(ctx context.Context, ttl time.Duration) error {
	return nil
}

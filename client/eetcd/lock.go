package eetcd

import (
	"context"
	"time"

	"github.com/coreos/etcd/clientv3/concurrency"
)

// Mutex ...
type Mutex struct {
	s *concurrency.Session
	m *concurrency.Mutex
}

// NewMutex ...
func (client *Component) NewMutex(key string, opts ...concurrency.SessionOption) (mutex *Mutex, err error) {
	mutex = &Mutex{}
	// 默认session ttl = 60s
	mutex.s, err = concurrency.NewSession(client.Client, opts...)
	if err != nil {
		return
	}
	mutex.m = concurrency.NewMutex(mutex.s, key)
	return
}

// Lock ...
func (mutex *Mutex) Lock(timeout time.Duration) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return mutex.m.Lock(ctx)
}

// TryLock ...
func (mutex *Mutex) TryLock(timeout time.Duration) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return mutex.m.Lock(ctx)
}

// Unlock ...
func (mutex *Mutex) Unlock() (err error) {
	err = mutex.m.Unlock(context.TODO())
	if err != nil {
		return
	}
	return mutex.s.Close()
}

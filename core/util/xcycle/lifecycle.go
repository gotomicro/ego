package xcycle

import (
	"sync"
	"sync/atomic"
)

// Cycle ..
type Cycle struct {
	mu      *sync.Mutex
	wg      *sync.WaitGroup
	done    chan struct{}
	quit    chan error
	closing uint32
	waiting uint32
	cnt     uint64
}

// NewCycle new a cycle life
func NewCycle() *Cycle {
	return &Cycle{
		mu:      &sync.Mutex{},
		wg:      &sync.WaitGroup{},
		done:    make(chan struct{}),
		quit:    make(chan error),
		closing: 0,
		waiting: 0,
	}
}

// Run a new goroutine
func (c *Cycle) Run(fn func() error) {
	c.mu.Lock()
	// todo add check options panic before waiting
	defer c.mu.Unlock()
	c.wg.Add(1)
	c.cnt++
	go func(c *Cycle) {
		defer c.wg.Done()
		if err := fn(); err != nil {
			c.quit <- err
		}
	}(c)
}

// Done block and return a chan error
func (c *Cycle) Done() <-chan struct{} {
	if atomic.CompareAndSwapUint32(&c.waiting, 0, 1) {
		go func(c *Cycle) {
			c.mu.Lock()
			defer c.mu.Unlock()
			c.wg.Wait()
			close(c.done)
		}(c)
	}
	return c.done
}

// DoneAndClose ..
func (c *Cycle) DoneAndClose() {
	<-c.Done()
	c.Close()
}

// Close ..
func (c *Cycle) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if atomic.CompareAndSwapUint32(&c.closing, 0, 1) {
		close(c.quit)
	}
}

// Wait blocked for a life cycle
func (c *Cycle) Wait(hang bool) <-chan error {
	c.mu.Lock()
	// 说明没有用户使用，直接关闭
	if c.cnt == 0 && !hang {
		if atomic.CompareAndSwapUint32(&c.closing, 0, 1) {
			close(c.quit)
		}
	}
	c.mu.Unlock()
	return c.quit
}

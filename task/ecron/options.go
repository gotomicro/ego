package ecron

import (
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/gotomicro/ego/core/elog"
)

// Option 可选项
type Option func(c *Container)

// WithLocker 注入分布式locker
func WithLocker(locker Locker) Option {
	return func(c *Container) {
		c.config.locker = locker
	}
}

// WithChain ...
func WithChain(wrappers ...JobWrapper) Option {
	return func(c *Container) {
		if c.config.wrappers == nil {
			c.config.wrappers = []JobWrapper{}
		}
		c.config.wrappers = append(c.config.wrappers, wrappers...)
	}
}

// queueIfStillRunning serializes jobs, delaying subsequent runs until the
// previous one is complete. Jobs running after a delay of more than a minute
// have the delay logged at Info.
func queueIfStillRunning(logger *elog.Component) JobWrapper {
	return func(j Job) Job {
		var mu sync.Mutex
		return cron.FuncJob(func() {
			start := time.Now()
			mu.Lock()
			defer mu.Unlock()
			if dur := time.Since(start); dur > time.Minute {
				logger.Info("cron queue", elog.String("duration", dur.String()))
			}
			j.Run()
		})
	}
}

// skipIfStillRunning skips an invocation of the Job if a previous invocation is
// still running. It logs skips to the given logger at Info level.
func skipIfStillRunning(logger *elog.Component) JobWrapper {
	var ch = make(chan struct{}, 1)
	ch <- struct{}{}
	return func(j Job) Job {
		return cron.FuncJob(func() {
			select {
			case v := <-ch:
				j.Run()
				ch <- v
			default:
				logger.Info("cron skip")
			}
		})
	}
}

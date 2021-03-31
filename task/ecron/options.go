package ecron

import (
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/gotomicro/ego/core/elog"
)

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

// Option ...
type Option func(c *Container)

// WithLock 设置分布式锁. 当 Config.EnableDistributedTask = true 时, 本 Option 必须设置
func WithLock(lock Lock) Option {
	return func(c *Container) {
		c.config.lock = lock
	}
}

// WithWrappers 设置 JobWrapper
func WithWrappers(wrappers ...JobWrapper) Option {
	return func(c *Container) {
		if c.config.wrappers == nil {
			c.config.wrappers = []JobWrapper{}
		}
		c.config.wrappers = append(c.config.wrappers, wrappers...)
	}
}

//WithJob 指定Job
func WithJob(job FuncJob) Option {
	return func(c *Container) {
		c.config.job = job
	}
}

//WithSeconds 开启秒单位
func WithSeconds() Option {
	return func(c *Container) {
		c.config.EnableSeconds = true
	}
}

//WithParser 设置时间 parser
func WithParser(p cron.Parser) Option {
	return func(c *Container) {
		c.config.parser = p
	}
}

//WithLocation 设置时区
func WithLocation(loc *time.Location) Option {
	return func(c *Container) {
		c.config.loc = loc
	}
}

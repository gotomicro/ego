package ecron

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/metric"
	"github.com/gotomicro/ego/core/util/xtime"
)

// Config ...
type Config struct {
	WithSeconds     bool
	ConcurrentDelay int
	ImmediatelyRun  bool
	TTL             int // 单位：s
	DistributedTask bool
	WaitLockTime    time.Duration
	WaitUnlockTime  time.Duration
	WorkerLockDir   string
	wrappers        []JobWrapper
	parser          cron.Parser
	locker          Locker
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		WithSeconds:     false,
		ConcurrentDelay: -1, // skip
		ImmediatelyRun:  false,
		TTL:             60,
		DistributedTask: false,
		WaitLockTime:    xtime.Duration("1000ms"),
		WaitUnlockTime:  xtime.Duration("1000ms"),
		WorkerLockDir:   "/ecron/lock/",
		wrappers:        []JobWrapper{},
		parser:          cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor),
	}
}

type wrappedLogger struct {
	*elog.Component
}

// Info logs routine messages about cron's operation.
func (wl *wrappedLogger) Info(msg string, keysAndValues ...interface{}) {
	wl.Infow("cron "+msg, keysAndValues...)
}

// Error logs an error condition.
func (wl *wrappedLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	wl.Errorw("cron "+msg, append(keysAndValues, "err", err)...)
}

type Locker interface {
	Lock(context.Context) error
	Unlock(context.Context) error
}

type wrappedJob struct {
	NamedJob
	logger          *elog.Component
	workerLockDir   string
	distributedTask bool
	waitLockTime    time.Duration
	waitUnlockTime  time.Duration
	leaseTTL        int
	locker          Locker
}

// Run ...
func (wj wrappedJob) Run() {
	if wj.distributedTask {
		var err error
		// 阻塞等待直到waitLockTime timeout
		ctx, cancel := context.WithTimeout(context.Background(), wj.waitLockTime)
		defer cancel()
		err = wj.locker.Lock(ctx)
		if err != nil {
			wj.logger.Info("mutex lock", elog.String("err", err.Error()))
			return
		}

		ctx, cancel = context.WithTimeout(context.Background(), wj.waitUnlockTime)
		defer cancel()
		err = wj.locker.Unlock(ctx)
		if err != nil {
			wj.logger.Info("mutex unlock", elog.String("err", err.Error()))
			return
		}
	}
	_ = wj.run()
}

func (wj wrappedJob) run() (err error) {
	metric.JobHandleCounter.Inc("cron", wj.Name(), "begin")
	var fields = []elog.Field{zap.String("name", wj.Name())}
	var beg = time.Now()
	defer func() {
		if rec := recover(); rec != nil {
			switch rec := rec.(type) {
			case error:
				err = rec
			default:
				err = fmt.Errorf("%v", rec)
			}

			stack := make([]byte, 4096)
			length := runtime.Stack(stack, true)
			fields = append(fields, zap.ByteString("stack", stack[:length]))
		}
		if err != nil {
			fields = append(fields, elog.String("err", err.Error()), elog.Duration("cost", time.Since(beg)))
			wj.logger.Error("run", fields...)
		} else {
			wj.logger.Info("run", fields...)
		}
		metric.JobHandleHistogram.Observe(time.Since(beg).Seconds(), "cron", wj.Name())
	}()

	return wj.NamedJob.Run()
}

package ecron

import (
	"fmt"
	"github.com/gotomicro/ego/core/util/xtime"
	"runtime"
	"time"

	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/gotomicro/ego/client/eetcd"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/metric"
	"github.com/robfig/cron/v3"

	"go.uber.org/zap"
)

// Config ...
type Config struct {
	WithSeconds         bool
	ConcurrentDelay     int
	ImmediatelyRun      bool
	TTL                 int // 单位：s
	DistributedTask     bool
	WaitLockTime        time.Duration
	WorkerLockDir       string
	DefaultWaitLockTime time.Duration //
	wrappers            []JobWrapper
	parser              cron.Parser
	// Distributed task
	etcdClient *eetcd.Component
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		WithSeconds:         false,
		ConcurrentDelay:     -1, // skip
		ImmediatelyRun:      false,
		TTL:                 60,
		DistributedTask:     false,
		WaitLockTime:        0,
		WorkerLockDir:       "/ecron/lock/",
		DefaultWaitLockTime: xtime.Duration("1000ms"), //ms
		wrappers:            []JobWrapper{},
		parser:              cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor),
		etcdClient:          nil,
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

type wrappedJob struct {
	NamedJob
	logger              *elog.Component
	workerLockDir       string
	distributedTask     bool
	waitLockTime        time.Duration
	leaseTTL            int
	client              *eetcd.Component
	defaultWaitLockTime time.Duration
}

// Run ...
func (wj wrappedJob) Run() {
	if wj.distributedTask {
		mutex, err := wj.client.NewMutex(wj.workerLockDir+wj.Name(), concurrency.WithTTL(wj.leaseTTL))
		if err != nil {
			wj.logger.Error("mutex", elog.String("err", err.Error()))
			return
		}
		if wj.waitLockTime == 0 {
			err = mutex.TryLock(wj.defaultWaitLockTime)
		} else { // 阻塞等待直到waitLockTime timeout
			err = mutex.Lock(wj.waitLockTime)
		}
		if err != nil {
			wj.logger.Info("mutex lock", elog.String("err", err.Error()))
			return
		}
		defer mutex.Unlock()
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

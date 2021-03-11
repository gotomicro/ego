package ecron

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/emetric"
	"github.com/gotomicro/ego/core/util/xtime"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Config ...
type Config struct {
	WaitLockTime          time.Duration // 抢锁等待时间，默认60s
	LockTTL               time.Duration // 租期，默认60s
	LockDir               string        // 定时任务锁目录
	RefreshTTL            time.Duration // 刷新ttl，默认60s
	WaitUnlockTime        time.Duration // 抢锁等待时间，默认1s
	DelayExecType         string        // skip，queue，concurrent，如果上一个任务执行较慢，到达了新任务执行时间，那么新任务选择跳过，排队，并发执行的策略，新任务默认选择skip策略
	EnableDistributedTask bool          // 是否分布式任务，默认否，如果存在分布式任务，会只执行该定时人物
	EnableImmediatelyRun  bool          // 是否立刻执行，默认否
	EnableWithSeconds     bool          // 是否使用秒作解析器，默认否
	wrappers              []JobWrapper
	parser                cron.Parser
	locker                Locker
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		EnableWithSeconds:     false,
		EnableImmediatelyRun:  false,
		EnableDistributedTask: false,
		DelayExecType:         "skip",
		WaitLockTime:          xtime.Duration("60s"),
		LockTTL:               xtime.Duration("60s"),
		RefreshTTL:            xtime.Duration("50s"),
		WaitUnlockTime:        xtime.Duration("1s"),
		LockDir:               "/ecron/lock/%s/%s",
		wrappers:              []JobWrapper{},
		parser:                cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor),
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

// go get github.com/gotomicro/eredis@v0.2.0+
type Locker interface {
	Lock(ctx context.Context, key string, ttl time.Duration) error
	Unlock(ctx context.Context, key string) error
	Refresh(ctx context.Context, key string, ttl time.Duration) error
}

type wrappedJob struct {
	NamedJob
	logger *elog.Component
}

// Run ...
func (wj wrappedJob) Run() {
	_ = wj.run()
	return
}

func (wj wrappedJob) run() (err error) {
	emetric.JobHandleCounter.Inc("cron", wj.Name(), "begin")
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
		emetric.JobHandleHistogram.Observe(time.Since(beg).Seconds(), "cron", wj.Name())
	}()

	return wj.NamedJob.Run()
}

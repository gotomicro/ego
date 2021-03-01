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
	WithSeconds      bool          // 是否使用秒作解析器，默认否
	ConcurrentDelay  int           // 并发延迟，默认是执行超过定时时间后，下次执行的任务会跳过
	ImmediatelyRun   bool          // 是否立刻执行，默认否
	DistributedTask  bool          // 是否分布式任务，默认否，如果存在分布式任务，则会解析嵌入的etcd配置
	WaitLockTime     time.Duration // 抢锁等待时间，默认60s
	LockTTL          time.Duration // 租期，默认60s
	LockDir          string        // 定时任务锁目录
	RefreshTTL       time.Duration // 刷新ttl，默认60s
	WaitUnlockTime   time.Duration // 抢锁等待时间，默认1s
	Endpoints        []string      // etcd地址
	ConnectTimeout   time.Duration // 连接超时时间，默认5s
	Secure           bool          // 是否安全通信，默认false
	AutoSyncInterval time.Duration // 自动同步member list的间隔
	wrappers         []JobWrapper
	parser           cron.Parser
	locker           Locker
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		WithSeconds:     false,
		ConcurrentDelay: -1, // skip
		ImmediatelyRun:  false,
		DistributedTask: false,
		WaitLockTime:    xtime.Duration("60s"),
		LockTTL:         xtime.Duration("60s"),
		RefreshTTL:      xtime.Duration("50s"),
		WaitUnlockTime:  xtime.Duration("1s"),
		LockDir:         "/ecron/lock/%s/%s",
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

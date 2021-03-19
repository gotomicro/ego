package ecron

import (
	"time"

	"github.com/robfig/cron/v3"

	"github.com/gotomicro/ego/core/util/xtime"
)

// Config ...
type Config struct {
	// Required. 触发时间
	//	默认最小单位为分钟.比如:
	//		"* * * * * *" 代表每分钟执行
	//	如果 EnableSeconds = true. 那么最小单位为秒. 示例:
	//		"*/3 * * * * * *" 代表每三秒钟执行一次
	Spec string

	WaitLockTime   time.Duration // 抢锁等待时间，默认 4s
	LockTTL        time.Duration // 租期，默认 16s
	RefreshGap     time.Duration // 锁刷新间隔时间， 默认 4s
	WaitUnlockTime time.Duration // 解锁等待时间，默认 1s

	DelayExecType         string // skip，queue，concurrent，如果上一个任务执行较慢，到达了新任务执行时间，那么新任务选择跳过，排队，并发执行的策略，新任务默认选择skip策略
	EnableDistributedTask bool   // 是否分布式任务，默认否，如果存在分布式任务，会只执行该定时人物
	EnableImmediatelyRun  bool   // 是否立刻执行，默认否
	EnableSeconds         bool   // 是否使用秒作解析器，默认否

	wrappers []JobWrapper
	parser   cron.Parser
	lock     Lock
	job      FuncJob
	loc      *time.Location
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Spec:                  "", // required in config
		WaitLockTime:          xtime.Duration("4s"),
		LockTTL:               xtime.Duration("16s"),
		RefreshGap:            xtime.Duration("4s"),
		WaitUnlockTime:        xtime.Duration("1s"),
		DelayExecType:         "skip",
		EnableDistributedTask: false,
		EnableImmediatelyRun:  false,
		EnableSeconds:         false,
		wrappers:              []JobWrapper{},
		parser:                cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor),
		lock:                  nil,
		job:                   nil,
		loc:                   time.Local,
	}
}

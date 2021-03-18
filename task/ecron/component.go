package ecron

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/standard"

	"github.com/robfig/cron/v3"

	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/util/xstring"
)

const PackageName = "task.ecron"

var (
	// Every ...
	Every = cron.Every
	// NewParser ...
	NewParser = cron.NewParser
	// NewChain ...
	NewChain = cron.NewChain
	// WithSeconds ...
	WithSeconds = cron.WithSeconds
	// WithParser ...
	WithParser = cron.WithParser
	// WithLocation ...
	WithLocation = cron.WithLocation
)

type (
	// JobWrapper ...
	JobWrapper = cron.JobWrapper
	// EntryID ...
	EntryID = cron.EntryID
	// Entry ...
	Entry = cron.Entry
	// Schedule ...
	Schedule = cron.Schedule
	// Parser ...
	Parser = cron.Parser

	// Job ...
	Job = cron.Job
	//NamedJob ..
	NamedJob interface {
		Run(ctx context.Context) error
		Name() string
	}
)

// FuncJob ...
type FuncJob func(ctx context.Context) error

// Run ...
func (f FuncJob) Run(ctx context.Context) error {
	return f(ctx)
}

// Name ...
func (f FuncJob) Name() string { return xstring.FunctionName(f) }

// Component ...
type Component struct {
	name   string
	config *Config
	*cron.Cron
	entries map[string]EntryID
	logger  *elog.Component
}

func newComponent(name string, config *Config, logger *elog.Component) *Component {
	cron := &Component{
		config: config,
		Cron: cron.New(
			cron.WithParser(config.parser),
			cron.WithChain(config.wrappers...),
			cron.WithLogger(&wrappedLogger{logger}),
		),
		name:   name,
		logger: logger,
	}
	return cron
}

// Schedule ...
func (c *Component) Schedule(schedule Schedule, job NamedJob) EntryID {
	if c.config.EnableImmediatelyRun {
		schedule = &immediatelyScheduler{
			Schedule: schedule,
		}
	}
	innnerJob := &wrappedJob{
		NamedJob: job,
		logger:   c.logger,
	}
	c.logger.Info("add job", elog.String("name", job.Name()))
	return c.Cron.Schedule(schedule, innnerJob)
}

func (c *Component) Name() string {
	return c.name
}

func (c *Component) PackageName() string {
	return PackageName
}

func (c *Component) Init() error {
	return nil
}

// GetEntryByName ...
func (c *Component) GetEntryByName(name string) cron.Entry {
	return c.Entry(c.entries[name])
}

// AddJob ...
func (c *Component) AddJob(spec string, cmd NamedJob) (EntryID, error) {
	schedule, err := c.config.parser.Parse(spec)
	if err != nil {
		return 0, err
	}
	return c.Schedule(schedule, cmd), nil
}

// AddFunc ...
func (c *Component) AddFunc(spec string, cmd func(ctx context.Context) error) (EntryID, error) {
	return c.AddJob(spec, FuncJob(cmd))
}

// Start ...
func (c *Component) Start() error {
	if c.config.EnableDistributedTask {
		// 如果分布式的定时任务，那么就需要抢占锁
		go func() {
			var err error
			for {
				// 阻塞等待直到waitLockTime timeout
				ctx, cancel := context.WithTimeout(context.Background(), c.config.WaitLockTime)
				defer cancel()
				err = c.config.locker.Lock(ctx, c.lockerName(), c.config.LockTTL)
				if err != nil {
					c.logger.Info("mutex lock", elog.FieldErr(err))
					continue
				}

				c.logger.Info("add cron", elog.Int("number of scheduled jobs", len(c.Cron.Entries())))

				// 定时续期
				go func() {
					for {
						// todo 网络异常情况导致的问题
						c.config.locker.Refresh(context.Background(), c.lockerName(), c.config.LockTTL)
						time.Sleep(c.config.RefreshTTL)
					}
				}()
				c.Cron.Run()
				return
			}
		}()

	} else {
		c.logger.Info("add cron", elog.Int("number of scheduled jobs", len(c.Cron.Entries())))
		c.Cron.Run()
	}

	return nil
}

// Stop ...
func (c *Component) Stop() error {
	_ = c.Cron.Stop()
	if c.config.EnableDistributedTask {
		ctx, cancel := context.WithTimeout(context.Background(), c.config.WaitUnlockTime)
		defer cancel()
		err := c.config.locker.Unlock(ctx, c.lockerName())
		if err != nil {
			c.logger.Info("mutex unlock", elog.FieldErr(err))
			return fmt.Errorf("cron stop err: %w", err)
		}
	}
	return nil
}

// locker name
//
func (c *Component) lockerName() string {
	return fmt.Sprintf(c.config.LockDir, eapp.AppInstance(), c.name)
}

type immediatelyScheduler struct {
	Schedule
	initOnce uint32
}

// Next ...
func (is *immediatelyScheduler) Next(curr time.Time) (next time.Time) {
	if atomic.CompareAndSwapUint32(&is.initOnce, 0, 1) {
		return curr
	}

	return is.Schedule.Next(curr)
}

type Ecron interface {
	standard.Component
}

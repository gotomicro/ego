package ecron

import (
	"context"
	"fmt"
	"time"

	"github.com/gotomicro/ego/core/etrace"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/util/xstring"
)

// PackageName 包名
const PackageName = "task.ecron"

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
	cron   *cron.Cron
	logger *elog.Component
}

func newComponent(name string, config *Config, logger *elog.Component) *Component {
	return &Component{
		config: config,
		cron: cron.New(
			cron.WithParser(config.parser),
			cron.WithChain(config.wrappers...),
			cron.WithLogger(&wrappedLogger{logger}),
			cron.WithLocation(config.loc),
		),
		name:   name,
		logger: logger,
	}
}

// Name 名称
func (c *Component) Name() string {
	return c.name
}

// PackageName 包名
func (c *Component) PackageName() string {
	return PackageName
}

// Init Init
func (c *Component) Init() error {
	return nil
}

// Start ...
func (c *Component) Start() error {
	if !c.config.Enable {
		return nil
	}

	if c.config.EnableDistributedTask {
		go c.startDistributedTask()
	} else {
		err := c.startTask()
		if err != nil {
			return err
		}
	}

	c.cron.Run()
	return nil
}

// Stop ...
func (c *Component) Stop() error {
	_ = c.cron.Stop()
	if c.config.EnableDistributedTask {
		ctx, cancel := context.WithTimeout(context.Background(), c.config.WaitUnlockTime)
		defer cancel()
		err := c.config.lock.Unlock(ctx)
		if err != nil {
			c.logger.Info("mutex unlock", elog.FieldErr(err))
			return fmt.Errorf("cron stop err: %w", err)
		}
	}
	return nil
}

func (c *Component) schedule(schedule Schedule, job NamedJob) EntryID {
	if c.config.EnableImmediatelyRun {
		schedule = &immediatelyScheduler{
			Schedule: schedule,
		}
	}
	innerJob := &wrappedJob{
		NamedJob: job,
		logger:   c.logger,
		tracer:   etrace.NewTracer(trace.SpanKindServer),
	}
	c.logger.Info("add job", elog.String("name", job.Name()))
	return c.cron.Schedule(schedule, innerJob)
}

func (c *Component) addJob(spec string, cmd NamedJob) (EntryID, error) {
	schedule, err := c.config.parser.Parse(spec)
	if err != nil {
		return 0, err
	}
	return c.schedule(schedule, cmd), nil
}

func (c *Component) removeJob(id EntryID) {
	c.cron.Remove(id)
}

func (c *Component) startDistributedTask() {
	for {
		func() {
			defer time.Sleep(c.config.RefreshGap)

			ctx, cancel := context.WithTimeout(context.Background(), c.config.WaitLockTime)
			err := c.config.lock.Lock(ctx, c.config.LockTTL)
			cancel()
			if err != nil {
				c.logger.Info("job lock not obtained", elog.FieldErr(err))
				return
			}

			c.logger.Info("add cron", elog.Int("number of scheduled jobs", len(c.cron.Entries())))

			entryID, err := c.addJob(c.config.Spec, c.config.job)
			if err != nil {
				c.logger.Error("add job failed", zap.Error(err))
				return
			}

			err = c.keepLockAlive()
			if err != nil {
				c.logger.Error("job lost", zap.String("name", c.name), zap.Error(err))
			}

			c.removeJob(entryID)
		}()
	}
}

func (c *Component) keepLockAlive() error {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), c.config.WaitLockTime)
		err := c.config.lock.Refresh(ctx, c.config.LockTTL)
		cancel()
		if err != nil {
			c.logger.Info("mutex lock", elog.FieldErr(err))
			return err
		}

		time.Sleep(c.config.RefreshGap)
	}
}

func (c *Component) startTask() (err error) {
	_, err = c.addJob(c.config.Spec, c.config.job)
	if err != nil {
		return
	}

	c.logger.Info("add cron", elog.Int("number of scheduled jobs", len(c.cron.Entries())))
	return nil
}

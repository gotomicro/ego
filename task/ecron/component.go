package ecron

import (
	"sync/atomic"
	"time"

	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/util/xstring"
	"github.com/robfig/cron/v3"
)

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
		Run() error
		Name() string
	}
)

// FuncJob ...
type FuncJob func() error

// Run ...
func (f FuncJob) Run() error { return f() }

// Name ...
func (f FuncJob) Name() string { return xstring.FunctionName(f) }

// Component ...
type Component struct {
	name   string
	Config *Config
	*cron.Cron
	entries map[string]EntryID
	logger  *elog.Component
}

func newComponent(name string, config *Config, logger *elog.Component) *Component {
	cron := &Component{
		Config: config,
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
	if c.Config.ImmediatelyRun {
		schedule = &immediatelyScheduler{
			Schedule: schedule,
		}
	}
	innnerJob := &wrappedJob{
		NamedJob:            job,
		logger:              c.logger,
		workerLockDir:       c.Config.WorkerLockDir,
		distributedTask:     c.Config.DistributedTask,
		waitLockTime:        c.Config.WaitLockTime,
		leaseTTL:            c.Config.TTL,
		client:              c.Config.etcdClient,
		defaultWaitLockTime: c.Config.DefaultWaitLockTime,
	}
	c.logger.Info("add job", elog.String("name", job.Name()))
	return c.Cron.Schedule(schedule, innnerJob)
}

// GetEntryByName ...
func (c *Component) GetEntryByName(name string) cron.Entry {
	// todo(gorexlv): data race
	return c.Entry(c.entries[name])
}

// AddJob ...
func (c *Component) AddJob(spec string, cmd NamedJob) (EntryID, error) {
	schedule, err := c.Config.parser.Parse(spec)
	if err != nil {
		return 0, err
	}
	return c.Schedule(schedule, cmd), nil
}

// AddFunc ...
func (c *Component) AddFunc(spec string, cmd func() error) (EntryID, error) {
	return c.AddJob(spec, FuncJob(cmd))
}

// Run ...
func (c *Component) Run() error {
	// xdebug.PrintKVWithPrefix("worker", "run worker", fmt.Sprintf("%d job scheduled", len(c.Component.Entries())))
	c.logger.Info("run worker", elog.Int("number of scheduled jobs", len(c.Cron.Entries())))
	c.Cron.Run()
	return nil
}

// Stop ...
func (c *Component) Stop() error {
	_ = c.Cron.Stop()
	return nil
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

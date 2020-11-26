package ego

import (
	"context"
	"github.com/gotomicro/ego/core/ecode"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/flag"
	"github.com/gotomicro/ego/core/registry"
	"github.com/gotomicro/ego/core/util/xcycle"
	"github.com/gotomicro/ego/server"
	"github.com/gotomicro/ego/task/ecron"
	"github.com/gotomicro/ego/task/ejob"
	"strings"
	"sync"
)

// Ego分为三大部分
// 第一部分 系统数据：生命周期，配置前缀，锁，日志，错误
// 第二部分 运行程序：系统初始化函数，用户初始化函数，服务，定时任务，短时任务
// 第三部分 可选方法：是否悬挂，注册中心，运行停止前清理，运行停止后清理
type ego struct {
	// 第一部分 系统数据
	cycle        *xcycle.Cycle   // 生命周期
	configPrefix string          // 配置前缀
	smu          *sync.RWMutex   // 锁
	logger       *elog.Component // 日志
	err          error           // 错误

	// 第二部分 运行程序
	inits    []func() error         // 系统初始化函数
	invokers []func() error         // 用户初始化函数
	servers  []server.Server        // 服务
	crons    []ecron.Cron           // 定时任务
	jobs     map[string]ejob.Runner // 短时任务

	// 第三部分 可选方法
	hang            bool              // 是否悬挂
	registerer      registry.Registry // 注册中心
	beforeStopClean []func() error    // 运行停止前清理
	afterStopClean  []func() error    // 运行停止后清理
}

// New new ego
func New(options ...Option) *ego {
	e := &ego{
		// 第一部分 系统数据
		cycle:        xcycle.NewCycle(),
		smu:          &sync.RWMutex{},
		configPrefix: "",
		logger:       elog.EgoLogger,
		err:          nil,

		// 第二部分 运行程序
		inits:    make([]func() error, 0),
		invokers: make([]func() error, 0),
		servers:  make([]server.Server, 0),
		crons:    make([]ecron.Cron, 0),
		jobs:     make(map[string]ejob.Runner),

		// 第三部分 可选方法
		hang:            false,
		registerer:      registry.Nop{},
		beforeStopClean: make([]func() error, 0),
		afterStopClean:  make([]func() error, 0),
	}

	// 设置运行前清理函数
	// 如果注册中心存在设置
	if e.registerer != nil {
		WithBeforeStopClean(e.registerer.Close)
	}

	// 设置运行后清理函数
	// 设置清理日志函数
	WithAfterStopClean(elog.DefaultLogger.Flush, elog.EgoLogger.Flush)

	// 设置参数
	for _, option := range options {
		option(e)
	}

	// 设置初始函数
	e.inits = []func() error{
		parseFlags,
		printBanner,
		loadConfig,
		initMaxProcs,
		e.initLogger,
		e.initTracer,
	}

	// 初始化系统函数
	e.err = runSerialFuncReturnError(e.inits)
	return e
}

// Invoker 传入所需要的函数
func (e *ego) Invoker(fns ...func() error) *ego {
	e.smu.Lock()
	defer e.smu.Unlock()

	e.invokers = append(e.invokers, fns...)

	// 初始化用户函数
	e.err = runSerialFuncReturnError(e.invokers)
	return e
}

// 服务
func (e *ego) Serve(s ...server.Server) *ego {
	e.smu.Lock()
	defer e.smu.Unlock()
	e.servers = append(e.servers, s...)
	return e
}

// 定时任务
func (e *ego) Cron(w ...ecron.Cron) *ego {
	e.crons = append(e.crons, w...)
	return e
}

// 短时任务
func (e *ego) Job(runners ...ejob.Runner) *ego {
	// start job by name
	jobFlag := flag.String("job")
	if jobFlag == "" {
		e.logger.Info("ego jobs flag name empty")
		return e
	}

	jobMap := make(map[string]struct{}, 0)
	// 逗号分割可以执行多个job
	if strings.Contains(jobFlag, ",") {
		jobArr := strings.Split(jobFlag, ",")
		for _, value := range jobArr {
			jobMap[value] = struct{}{}
		}
	} else {
		jobMap[jobFlag] = struct{}{}
	}

	for _, runner := range runners {
		namedJob, ok := runner.(interface{ GetJobName() string })
		// job runner must implement GetJobName
		if !ok {
			return e
		}
		jobName := namedJob.GetJobName()
		if flag.Bool("disable-job") {
			e.logger.Info("ego disable job", elog.FieldName(jobName))
			return e
		}

		_, flag := jobMap[jobName]
		if flag {
			e.logger.Info("ego register job", elog.FieldName(jobName))
			e.jobs[jobName] = runner
		}
	}
	return e
}

// 运行程序
func (e *ego) Run() error {
	if e.err != nil {
		return e.err
	}

	// 如果存在短时任务，那么只执行短时任务
	if len(e.jobs) > 0 {
		return e.startJobs()
	}

	e.waitSignals() // start signal listen task in goroutine

	// 启动服务
	e.startServers()

	// 启动定时任务
	e.startCrons()

	// 阻塞，等待信号量
	if err := <-e.cycle.Wait(e.hang); err != nil {
		e.logger.Error("ego shutdown with error", elog.FieldMod(ecode.ModApp), elog.FieldErr(err))
		return err
	}
	e.logger.Info("shutdown ego, bye!", elog.FieldMod(ecode.ModApp))

	// 运行停止后清理
	runSerialFuncLogError(e.afterStopClean)
	return nil
}

// 停止程序
func (e *ego) Stop(ctx context.Context, isGraceful bool) (err error) {
	// 运行停止前清理
	runSerialFuncLogError(e.beforeStopClean)

	// 停止服务
	e.smu.RLock()
	if isGraceful {
		for _, s := range e.servers {
			func(s server.Server) {
				e.cycle.Run(func() error {
					return s.GracefulStop(ctx)
				})
			}(s)
		}
	} else {
		for _, s := range e.servers {
			func(s server.Server) {
				e.cycle.Run(s.Stop)
			}(s)
		}
	}

	e.smu.RUnlock()

	// 停止定时任务
	for _, w := range e.crons {
		func(w ecron.Cron) {
			e.cycle.Run(w.Stop)
		}(w)
	}
	<-e.cycle.Done()
	e.cycle.Close()
	return err
}

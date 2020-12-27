package ego

import (
	"context"
	"github.com/gotomicro/ego/core/eflag"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/eregistry"
	"github.com/gotomicro/ego/core/util/xcycle"
	"github.com/gotomicro/ego/core/util/xtime"
	"github.com/gotomicro/ego/server"
	"github.com/gotomicro/ego/task/ecron"
	"github.com/gotomicro/ego/task/ejob"
	"os"
	"strings"
	"sync"
	"time"
)

// Ego分为三大部分
// 第一部分 系统数据：生命周期，配置前缀，锁，日志，错误
// 第二部分 运行程序：系统初始化函数，用户初始化函数，服务，定时任务，短时任务
// 第三部分 可选方法：是否悬挂，注册中心，运行停止前清理，运行停止后清理
type ego struct {
	// 第一部分 系统数据
	cycle  *xcycle.Cycle   // 生命周期
	smu    *sync.RWMutex   // 锁
	logger *elog.Component // 日志
	err    error           // 错误

	// 第二部分 运行程序
	inits      []func() error       // 系统初始化函数
	invokers   []func() error       // 用户初始化函数
	servers    []server.Server      // 服务
	crons      []ecron.Ecron        // 定时任务
	jobs       map[string]ejob.Ejob // 短时任务
	registerer eregistry.Registry   // 注册中心

	// 第三部分 可选方法
	opts opts
}

type opts struct {
	configPrefix    string         // 配置前缀
	hang            bool           // 是否悬挂
	disableBanner   bool           // 禁用banner
	beforeStopClean []func() error // 运行停止前清理
	afterStopClean  []func() error // 运行停止后清理
	stopTimeout     time.Duration  // 运行停止超时时间
	shutdownSignals []os.Signal
}

// New new ego
func New(options ...Option) *ego {
	e := &ego{
		// 第一部分 系统数据
		cycle:  xcycle.NewCycle(),
		smu:    &sync.RWMutex{},
		logger: elog.EgoLogger,
		err:    nil,

		// 第二部分 运行程序
		inits:      make([]func() error, 0),
		invokers:   make([]func() error, 0),
		servers:    make([]server.Server, 0),
		crons:      make([]ecron.Ecron, 0),
		jobs:       make(map[string]ejob.Ejob),
		registerer: eregistry.Nop{},

		// 第三部分 可选方法
		opts: opts{
			hang:            false,
			configPrefix:    "",
			beforeStopClean: make([]func() error, 0),
			afterStopClean:  make([]func() error, 0),
			stopTimeout:     xtime.Duration("5s"),
			shutdownSignals: shutdownSignals,
		},
	}

	// 设置运行前清理函数
	// 如果注册中心存在设置
	if e.registerer != nil {
		options = append(options, WithBeforeStopClean(e.registerer.Close))
	}

	// 设置运行后清理函数
	// 设置清理日志函数
	options = append(options, WithAfterStopClean(elog.DefaultLogger.Flush, elog.EgoLogger.Flush))

	// 设置参数
	for _, option := range options {
		option(e)
	}

	// 设置初始函数
	e.inits = []func() error{
		parseFlags,
		e.printBanner,
		printLogger,
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

func (e *ego) Registry(reg eregistry.Registry) *ego {
	e.registerer = reg
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
func (e *ego) Cron(w ...ecron.Ecron) *ego {
	e.crons = append(e.crons, w...)
	return e
}

// 短时任务
func (e *ego) Job(runners ...ejob.Ejob) *ego {
	// start job by name
	jobFlag := eflag.String("job")
	if jobFlag == "" {
		e.logger.Info("flag jobs name empty", elog.FieldComponent(ejob.PackageName))
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
		jobName := runner.Name()
		if jobName == "" {
			e.logger.Error("runner job name empty", elog.FieldComponent(runner.PackageName()))
			return e
		}
		if eflag.Bool("disable-job") {
			e.logger.Info("runner disable job", elog.FieldComponent(runner.PackageName()), elog.FieldName(jobName))
			return e
		}

		_, flag := jobMap[jobName]
		if flag {
			e.logger.Info("init register job", elog.FieldComponent(runner.PackageName()), elog.FieldName(jobName))
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
	if err := <-e.cycle.Wait(e.opts.hang); err != nil {
		e.logger.Error("ego shutdown with error", elog.FieldComponent("app"), elog.FieldErr(err))
		return err
	}
	e.logger.Info("stop ego, bye!", elog.FieldComponent("app"))

	// 运行停止后清理
	runSerialFuncLogError(e.opts.afterStopClean)
	return nil
}

// 停止程序
func (e *ego) Stop(ctx context.Context, isGraceful bool) (err error) {
	// 运行停止前清理
	runSerialFuncLogError(e.opts.beforeStopClean)

	// 停止服务
	e.smu.RLock()
	if isGraceful {
		for _, s := range e.servers {
			func(s server.Server) {
				// todo
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
		func(w ecron.Ecron) {
			e.cycle.Run(w.Stop)
		}(w)
	}
	<-e.cycle.Done()
	e.cycle.Close()
	return err
}

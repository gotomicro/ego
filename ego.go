// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ego

import (
	"context"
	"fmt"
	"github.com/gotomicro/ego/core/app"
	"github.com/gotomicro/ego/core/conf"
	"github.com/gotomicro/ego/core/conf/file"
	"github.com/gotomicro/ego/core/conf/manager"
	"github.com/gotomicro/ego/core/ecode"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/flag"
	"github.com/gotomicro/ego/core/registry"
	"github.com/gotomicro/ego/core/signals"
	"github.com/gotomicro/ego/core/trace"
	"github.com/gotomicro/ego/core/trace/jaeger"
	"github.com/gotomicro/ego/core/util/xcolor"
	"github.com/gotomicro/ego/core/util/xcycle"
	"github.com/gotomicro/ego/core/util/xdefer"
	"github.com/gotomicro/ego/core/util/xgo"
	"github.com/gotomicro/ego/server"
	"github.com/gotomicro/ego/task/ecron"
	job "github.com/gotomicro/ego/task/ejob"
	"go.uber.org/automaxprocs/maxprocs"
	"os"
	"runtime"
	"strings"
	"sync"
)

const (
	// StageAfterStop after app stop
	StageAfterStop uint32 = iota + 1
	// StageBeforeStop before app stop
	StageBeforeStop
)

// Application is the framework's instance, it contains the servers, crons, client and configuration settings.
// Create an instance of Application, by using &Application{}
type Application struct {
	cycle       *xcycle.Cycle
	smu         *sync.RWMutex
	initOnce    sync.Once
	startupOnce sync.Once
	stopOnce    sync.Once
	prefix      string // 配置前缀

	// 核心运行程序
	servers      []server.Server
	crons        []ecron.Cron
	jobs         map[string]job.Runner
	hang         bool
	logger       *elog.Component
	registerer   registry.Registry
	hooks        map[uint32]*xdefer.DeferStack
	configParser conf.Unmarshaller
	err          error
}

// New new a Application
func New(fns ...func() error) *Application {
	app := &Application{}
	return app.Invoker(fns...)
}

// init hooks
func (a *Application) initHooks(hookKeys ...uint32) {
	a.hooks = make(map[uint32]*xdefer.DeferStack, len(hookKeys))
	for _, k := range hookKeys {
		a.hooks[k] = xdefer.NewStack()
	}
}

// run hooks
func (a *Application) runHooks(k uint32) {
	hooks, ok := a.hooks[k]
	if ok {
		hooks.Clean()
	}
}

// RegisterHooks register a stage Hook
func (a *Application) RegisterHooks(k uint32, fns ...func() error) error {
	hooks, ok := a.hooks[k]
	if ok {
		hooks.Push(fns...)
		return nil
	}
	return fmt.Errorf("hook stage not found")
}

// initialize application
func (a *Application) initialize() {
	a.initOnce.Do(func() {
		// assign
		a.cycle = xcycle.NewCycle()
		a.smu = &sync.RWMutex{}
		a.servers = make([]server.Server, 0)
		a.crons = make([]ecron.Cron, 0)
		a.jobs = make(map[string]job.Runner)
		a.logger = elog.EgoLogger
		// private method
		a.initHooks(StageBeforeStop, StageAfterStop)
		// public method
		a.SetRegistry(registry.Nop{}) // default nop without registry
	})
}

// start up application
// By default the startup composition is:
// - parse config, watch, version flags
// - load config
// - init default biz logger, ego frame logger
// - init procs
func (a *Application) startup() (err error) {
	a.startupOnce.Do(func() {
		err = xgo.SerialUntilError(
			a.parseFlags,
			a.printBanner,
			a.loadConfig,
			a.initLogger,
			a.initMaxProcs,
			a.initTracer,
		)()
	})
	return
}

// Invoker 全局Error
func (a *Application) Invoker(fns ...func() error) *Application {
	a.initialize()
	if err := a.startup(); err != nil {
		a.err = err
		return a
	}
	a.err = xgo.SerialUntilError(fns...)()
	return a
}

// Serve start server
func (a *Application) Serve(s ...server.Server) *Application {
	a.smu.Lock()
	defer a.smu.Unlock()
	a.servers = append(a.servers, s...)
	return a
}

// Schedule ..
func (a *Application) Cron(w ...ecron.Cron) *Application {
	a.crons = append(a.crons, w...)
	return a
}

// Job ..
func (a *Application) Job(runners ...job.Runner) *Application {
	// start job by name
	jobFlag := flag.String("job")
	if jobFlag == "" {
		a.logger.Info("ego jobs flag name empty")
		return a
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
			return a
		}
		jobName := namedJob.GetJobName()
		if flag.Bool("disable-job") {
			a.logger.Info("ego disable job", elog.FieldName(jobName))
			return a
		}

		_, flag := jobMap[jobName]
		if flag {
			a.logger.Info("ego register job", elog.FieldName(jobName))
			a.jobs[jobName] = runner
		}
	}
	return a
}

// 是否允许系统悬挂起来，0 表示不悬挂， 1 表示悬挂。目的是一些脚本操作的时候，不想主线程停止
func (a *Application) Hang(flag bool) *Application {
	a.hang = flag
	return a
}

// SetRegistry set customize registry
func (a *Application) SetRegistry(reg registry.Registry) *Application {
	a.registerer = reg
	return a
}

// Run run application
func (a *Application) Run() error {
	if a.err != nil {
		return a.err
	}

	if len(a.jobs) > 0 {
		return a.startJobs()
	}

	a.waitSignals() // start signal listen task in goroutine
	defer a.clean()

	// start servers and govern server
	a.startServers()

	// start crons
	a.startWorkers()

	// blocking and wait quit
	if err := <-a.cycle.Wait(a.hang); err != nil {
		a.logger.Error("ego shutdown with error", elog.FieldMod(ecode.ModApp), elog.FieldErr(err))
		return err
	}
	a.logger.Info("shutdown ego, bye!", elog.FieldMod(ecode.ModApp))
	return nil
}

// clean after app quit
func (a *Application) clean() {
	_ = elog.DefaultLogger.Flush()
	_ = elog.EgoLogger.Flush()
}

// Stop application immediately after necessary cleanup
func (a *Application) Stop() (err error) {
	a.stopOnce.Do(func() {
		a.runHooks(StageBeforeStop)

		if a.registerer != nil {
			err = a.registerer.Close()
			if err != nil {
				a.logger.Error("stop register close err", elog.FieldMod(ecode.ModApp), elog.FieldErr(err))
			}
		}
		// stop servers
		a.smu.RLock()
		for _, s := range a.servers {
			func(s server.Server) {
				a.cycle.Run(s.Stop)
			}(s)
		}
		a.smu.RUnlock()

		// stop crons
		for _, w := range a.crons {
			func(w ecron.Cron) {
				a.cycle.Run(w.Stop)
			}(w)
		}
		<-a.cycle.Done()
		a.runHooks(StageAfterStop)
		a.cycle.Close()
	})
	return
}

// GracefulStop application after necessary cleanup
func (a *Application) GracefulStop(ctx context.Context) (err error) {
	a.stopOnce.Do(func() {
		a.runHooks(StageBeforeStop)

		if a.registerer != nil {
			err = a.registerer.Close()
			if err != nil {
				a.logger.Error("stop register close err", elog.FieldMod(ecode.ModApp), elog.FieldErr(err))
			}
		}
		// stop servers
		a.smu.RLock()
		for _, s := range a.servers {
			func(s server.Server) {
				a.cycle.Run(func() error {
					return s.GracefulStop(ctx)
				})
			}(s)
		}
		a.smu.RUnlock()

		// stop crons
		for _, w := range a.crons {
			func(w ecron.Cron) {
				a.cycle.Run(w.Stop)
			}(w)
		}
		<-a.cycle.Done()
		a.runHooks(StageAfterStop)
		a.cycle.Close()
	})
	return err
}

// waitSignals wait signal
func (a *Application) waitSignals() {
	a.logger.Info("init listen signal", elog.FieldMod(ecode.ModApp), elog.FieldEvent("init"))
	signals.Shutdown(func(grace bool) { // when get shutdown signal
		// todo: support timeout
		if grace {
			a.GracefulStop(context.TODO())
		} else {
			a.Stop()
		}
	})
}

func (a *Application) startServers() error {
	//var eg errgroup.Group
	// start multi servers
	for _, s := range a.servers {
		s := s
		a.cycle.Run(func() (err error) {
			s.Init()
			err = a.registerer.RegisterService(context.TODO(), s.Info())
			if err != nil {
				a.logger.Error("register service err", elog.FieldErr(err))
			}
			defer a.registerer.UnregisterService(context.TODO(), s.Info())
			a.logger.Info("start server", elog.FieldMod(ecode.ModApp), elog.FieldEvent("init"), elog.FieldName(s.Info().Name), elog.FieldAddr(s.Info().Label()), elog.Any("scheme", s.Info().Scheme))
			defer a.logger.Info("exit server", elog.FieldMod(ecode.ModApp), elog.FieldEvent("exit"), elog.FieldName(s.Info().Name), elog.FieldErr(err), elog.FieldAddr(s.Info().Label()))
			err = s.Start()
			return
		})
	}
	return nil
}

func (a *Application) startWorkers() error {
	// start multi crons
	for _, w := range a.crons {
		w := w
		a.cycle.Run(func() error {
			return w.Run()
		})
	}
	return nil
}

// todo handle error
func (a *Application) startJobs() error {
	if len(a.jobs) == 0 {
		return nil
	}
	var jobs = make([]func() error, 0)
	// warp jobs
	for name, runner := range a.jobs {
		jobs = append(jobs, func() error {
			a.logger.Info("job run begin", elog.FieldName(name))
			defer a.logger.Info("job run end", elog.FieldName(name))
			// runner.Run panic 错误在更上层抛出
			return runner.Run()
		})
	}

	return xgo.ParallelWithError(jobs...)()
}

// parseFlags init
func (a *Application) parseFlags() error {
	flag.Register(&flag.StringFlag{
		Name:    "config",
		Usage:   "--config",
		EnvVar:  "CONFIG",
		Default: "",
		Action:  func(name string, fs *flag.FlagSet) {},
	})

	flag.Register(&flag.BoolFlag{
		Name:    "watch",
		Usage:   "--watch, watch config change event",
		Default: true,
		EnvVar:  "CONFIG_WATCH",
	})

	flag.Register(&flag.BoolFlag{
		Name:    "version",
		Usage:   "--version, print version",
		Default: false,
		Action: func(string, *flag.FlagSet) {
			app.PrintVersion()
			os.Exit(0)
		},
	})

	flag.Register(&flag.StringFlag{
		Name:    "host",
		Usage:   "--host, print host",
		Default: "",
		Action:  func(string, *flag.FlagSet) {},
	})
	return flag.Parse()
}

// loadConfig init
func (a *Application) loadConfig() error {
	var configAddr = flag.String("config")
	// 如果配置为空，那么赋值默认配置
	if configAddr == "" {
		configAddr = app.EgoConfigPath()
	}

	// 暂时只支持文件
	file.Register()
	provider, err := manager.NewDataSource(file.DataSourceFile, configAddr, flag.Bool("watch"))
	if err != manager.ErrDefaultConfigNotExist {
		if err != nil {
			a.logger.Panic("data source: provider error", elog.FieldMod(ecode.ModConfig), elog.FieldErr(err))
		}

		parser, tag := file.ExtParser(configAddr)
		// 如果不是，就要加载文件，加载不到panic
		if err := conf.LoadFromDataSource(provider, parser, conf.TagName(tag)); err != nil {
			a.logger.Panic("data source: load config", elog.FieldMod(ecode.ModConfig), elog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), elog.FieldErr(err))
		}
		// 如果协议是file类型，并且是默认文件配置，那么判断下文件是否存在，如果不存在只告诉warning，什么都不做
	} else {
		a.logger.Info("no config... ", elog.FieldMod(ecode.ModConfig), elog.String("addr", configAddr))
	}
	return nil
}

// initLogger init
func (a *Application) initLogger() error {
	if conf.Get(a.prefix+"logger.default") != nil {
		elog.DefaultLogger = elog.Load(a.prefix + "logger.default").Build()
	}

	if conf.Get(a.prefix+"logger.ego") != nil {
		elog.EgoLogger = elog.Load(a.prefix + "logger.ego").Build(elog.WithFileName(elog.EgoLoggerName))
	}
	return nil
}

// initTracer init
func (a *Application) initTracer() error {
	// init tracing component jaeger
	if conf.Get(a.prefix+"trace.jaeger") != nil {
		var config = jaeger.RawConfig(a.prefix + "trace.jaeger")
		trace.SetGlobalTracer(config.Build())
	}
	return nil
}

// initMaxProcs init
func (a *Application) initMaxProcs() error {
	if maxProcs := conf.GetInt("maxProc"); maxProcs != 0 {
		runtime.GOMAXPROCS(maxProcs)
	} else {
		if _, err := maxprocs.Set(); err != nil {
			a.logger.Panic("auto max procs", elog.FieldMod(ecode.ModProc), elog.FieldErrKind(ecode.ErrKindAny), elog.FieldErr(err))
		}
	}
	a.logger.Info("auto max procs", elog.FieldMod(ecode.ModProc), elog.Int64("procs", int64(runtime.GOMAXPROCS(-1))))
	return nil
}

// printBanner init
func (a *Application) printBanner() error {
	const banner = `
 Welcome to Ego, starting application ...
`
	fmt.Println(xcolor.Green(banner))
	return nil
}

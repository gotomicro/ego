package ego

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.uber.org/automaxprocs/maxprocs"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/econf/manager"
	"github.com/gotomicro/ego/core/eflag"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/etrace"
	"github.com/gotomicro/ego/core/etrace/ejaeger"
	"github.com/gotomicro/ego/core/util/xcolor"
	"github.com/gotomicro/ego/core/util/xgo"
)

// waitSignals wait signal
func (e *Ego) waitSignals() {
	sig := make(chan os.Signal, 2)
	signal.Notify(
		sig,
		e.opts.shutdownSignals...,
	)
	go func() {
		s := <-sig
		grace := s != syscall.SIGQUIT
		go func() {
			stopCtx, cancel := context.WithTimeout(context.Background(), e.opts.stopTimeout)
			defer cancel()
			e.Stop(stopCtx, grace)
		}()
		<-sig
		os.Exit(128 + int(s.(syscall.Signal))) // second signal. Exit directly.
	}()
}

func (e *Ego) startServers() error {
	// start multi servers
	for _, s := range e.servers {
		s := s
		e.cycle.Run(func() (err error) {
			s.Init()
			err = e.registerer.RegisterService(context.TODO(), s.Info())
			if err != nil {
				e.logger.Error("register service err", elog.FieldComponent(s.PackageName()), elog.FieldComponentName(s.Name()), elog.FieldErr(err))
			}
			defer e.registerer.UnregisterService(context.TODO(), s.Info())
			e.logger.Info("start server", elog.FieldComponent(s.PackageName()), elog.FieldComponentName(s.Name()), elog.FieldAddr(s.Info().Label()))
			defer e.logger.Info("stop server", elog.FieldComponent(s.PackageName()), elog.FieldComponentName(s.Name()), elog.FieldErr(err), elog.FieldAddr(s.Info().Label()))
			err = s.Start()
			return
		})
	}
	return nil
}

func (e *Ego) startCrons() error {
	for _, w := range e.crons {
		w := w
		e.cycle.Run(func() error {
			return w.Start()
		})
	}
	return nil
}

// todo handle error
func (e *Ego) startJobs() error {
	if len(e.jobs) == 0 {
		return nil
	}
	var jobs = make([]func() error, 0)
	// warp jobs
	for _, runner := range e.jobs {
		runner := runner
		jobs = append(jobs, func() error {
			return runner.Start()
		})
	}
	return xgo.ParallelWithError(jobs...)()
}

// parseFlags init
func parseFlags() error {
	eflag.Register(&eflag.StringFlag{
		Name:    "config",
		Usage:   "--config",
		EnvVar:  "CONFIG",
		Default: "",
		Action:  func(name string, fs *eflag.FlagSet) {},
	})

	eflag.Register(&eflag.BoolFlag{
		Name:    "watch",
		Usage:   "--watch, watch config change event",
		Default: true,
		EnvVar:  "CONFIG_WATCH",
	})

	eflag.Register(&eflag.BoolFlag{
		Name:    "version",
		Usage:   "--version, print version",
		Default: false,
		Action: func(string, *eflag.FlagSet) {
			eapp.PrintVersion()
			os.Exit(0)
		},
	})

	eflag.Register(&eflag.StringFlag{
		Name:    "host",
		Usage:   "--host, print host",
		Default: "",
		Action:  func(string, *eflag.FlagSet) {},
	})
	return eflag.Parse()
}

// loadConfig init
func loadConfig() error {
	var configAddr = eflag.String("config")
	provider, parser, tag, err := manager.NewDataSource(configAddr, eflag.Bool("watch"))

	// 如果不存在配置，找不到该文件路径，该错误只存在file类型
	if err == manager.ErrDefaultConfigNotExist {
		// 如果协议是file类型，并且是默认文件配置，那么判断下文件是否存在，如果不存在只告诉warning，什么都不做
		elog.EgoLogger.Warn("no config... ", elog.FieldComponent(econf.PackageName), elog.String("addr", configAddr), elog.FieldErr(err))
		return nil
	}

	// 如果存在错误，报错
	if err != nil {
		elog.EgoLogger.Panic("data source: provider error", elog.FieldComponent(econf.PackageName), elog.FieldErr(err))
	}

	// 如果不是，就要加载文件，加载不到panic
	if err := econf.LoadFromDataSource(provider, parser, econf.WithTagName(tag)); err != nil {
		elog.EgoLogger.Panic("data source: load config", elog.FieldComponent(econf.PackageName), elog.FieldErrKind("unmarshal config err"), elog.FieldErr(err))
	}
	elog.EgoLogger.Info("init config", elog.FieldComponent(econf.PackageName), elog.String("addr", configAddr))
	return nil
}

// initLogger init application and Ego logger
func (e *Ego) initLogger() error {
	if econf.Get(e.opts.configPrefix+"logger.default") != nil {
		elog.DefaultLogger = elog.Load(e.opts.configPrefix + "logger.default").Build()
		elog.EgoLogger.Info("reinit default logger", elog.FieldComponent(elog.PackageName))
	}

	if econf.Get(e.opts.configPrefix+"logger.Ego") != nil {
		elog.EgoLogger = elog.Load(e.opts.configPrefix + "logger.Ego").Build(elog.WithFileName(elog.EgoLoggerName))
		elog.EgoLogger.Info("reinit Ego logger", elog.FieldComponent(elog.PackageName))
	}
	return nil
}

// initTracer init global tracer
func (e *Ego) initTracer() error {
	if econf.Get(e.opts.configPrefix+"trace.jaeger") != nil {
		container := ejaeger.Load(e.opts.configPrefix + "trace.jaeger")
		tracer := container.Build()
		etrace.SetGlobalTracer(tracer)
		e.opts.afterStopClean = append(e.opts.afterStopClean, container.Stop)
		elog.EgoLogger.Info("init trace", elog.FieldComponent("app"))
	} else {
		elog.EgoLogger.Warn("not init trace", elog.FieldValue("need trace config: [trace.jaeger] in config"), elog.FieldComponent("app"))
	}
	return nil
}

// initMaxProcs init
func initMaxProcs() error {
	if maxProcs := econf.GetInt("Ego.maxProc"); maxProcs != 0 {
		runtime.GOMAXPROCS(maxProcs)
	} else {
		if _, err := maxprocs.Set(); err != nil {
			elog.EgoLogger.Panic("init max procs", elog.FieldComponent("app"), elog.FieldErr(err))
		}
	}
	elog.EgoLogger.Info("init max procs", elog.FieldComponent("app"), elog.FieldValueAny(runtime.GOMAXPROCS(-1)))
	return nil
}

func printLogger() error {
	elog.EgoLogger.Info("init default logger", elog.FieldComponent(elog.PackageName))
	elog.EgoLogger.Info("init Ego logger", elog.FieldComponent(elog.PackageName))
	return nil
}

// printBanner init
func (e *Ego) printBanner() error {
	if e.opts.disableBanner {
		return nil
	}
	const banner = `
    _/_/_/_/    _/_/_/    _/_/   
   _/        _/        _/    _/   
  _/_/_/    _/  _/_/  _/    _/    
 _/        _/    _/  _/    _/     
_/_/_/_/    _/_/_/    _/_/  

Welcome to Ego, starting application ...
`
	fmt.Println(xcolor.Blue(banner))
	return nil
}

func runSerialFuncReturnError(fns []func() error) error {
	for _, fn := range fns {
		err := fn()
		if err != nil {
			return err
		}
	}
	return nil
}

func runSerialFuncLogError(fns []func() error) {
	for _, clean := range fns {
		err := clean()
		if err != nil {
			elog.EgoLogger.Error("beforeStopClean err", elog.FieldComponent("app"), elog.FieldErr(err))
		}
	}
}

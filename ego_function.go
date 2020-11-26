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
	"github.com/gotomicro/ego/core/signals"
	"github.com/gotomicro/ego/core/trace"
	"github.com/gotomicro/ego/core/trace/jaeger"
	"github.com/gotomicro/ego/core/util/xcolor"
	"github.com/gotomicro/ego/core/util/xgo"
	"go.uber.org/automaxprocs/maxprocs"
	"os"
	"runtime"
)

// waitSignals wait signal
func (e *ego) waitSignals() {
	e.logger.Info("init listen signal", elog.FieldMod(ecode.ModApp), elog.FieldEvent("init"))
	signals.Shutdown(func(grace bool) { // when get shutdown signal
		// todo: support timeout
		e.Stop(context.TODO(), grace)

	})
}

func (e *ego) startServers() error {
	//var eg errgroup.Group
	// start multi servers
	for _, s := range e.servers {
		s := s
		e.cycle.Run(func() (err error) {
			s.Init()
			err = e.registerer.RegisterService(context.TODO(), s.Info())
			if err != nil {
				e.logger.Error("register service err", elog.FieldErr(err))
			}
			defer e.registerer.UnregisterService(context.TODO(), s.Info())
			e.logger.Info("start server", elog.FieldMod(ecode.ModApp), elog.FieldEvent("init"), elog.FieldName(s.Info().Name), elog.FieldAddr(s.Info().Label()), elog.Any("scheme", s.Info().Scheme))
			defer e.logger.Info("exit server", elog.FieldMod(ecode.ModApp), elog.FieldEvent("exit"), elog.FieldName(s.Info().Name), elog.FieldErr(err), elog.FieldAddr(s.Info().Label()))
			err = s.Start()
			return
		})
	}
	return nil
}

func (e *ego) startCrons() error {
	// start multi crons
	for _, w := range e.crons {
		w := w
		e.cycle.Run(func() error {
			return w.Run()
		})
	}
	return nil
}

// todo handle error
func (e *ego) startJobs() error {
	if len(e.jobs) == 0 {
		return nil
	}
	var jobs = make([]func() error, 0)
	// warp jobs
	for name, runner := range e.jobs {
		jobs = append(jobs, func() error {
			e.logger.Info("job run begin", elog.FieldName(name))
			defer e.logger.Info("job run end", elog.FieldName(name))
			// runner.Run panic 错误在更上层抛出
			return runner.Run()
		})
	}

	return xgo.ParallelWithError(jobs...)()
}

// parseFlags init
func parseFlags() error {
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
func loadConfig() error {
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
			elog.EgoLogger.Panic("data source: provider error", elog.FieldMod(ecode.ModConfig), elog.FieldErr(err))
		}

		parser, tag := file.ExtParser(configAddr)
		// 如果不是，就要加载文件，加载不到panic
		if err := conf.LoadFromDataSource(provider, parser, conf.TagName(tag)); err != nil {
			elog.EgoLogger.Panic("data source: load config", elog.FieldMod(ecode.ModConfig), elog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), elog.FieldErr(err))
		}
		// 如果协议是file类型，并且是默认文件配置，那么判断下文件是否存在，如果不存在只告诉warning，什么都不做
	} else {
		elog.EgoLogger.Warn("no config... ", elog.FieldMod(ecode.ModConfig), elog.String("addr", configAddr))
	}
	return nil
}

// initLogger init
func (e *ego) initLogger() error {
	if conf.Get(e.configPrefix+"logger.default") != nil {
		elog.DefaultLogger = elog.Load(e.configPrefix + "logger.default").Build()
	}

	if conf.Get(e.configPrefix+"logger.ego") != nil {
		elog.EgoLogger = elog.Load(e.configPrefix + "logger.ego").Build(elog.WithFileName(elog.EgoLoggerName))
	}
	return nil
}

// initTracer init
func (e *ego) initTracer() error {
	// init tracing component jaeger
	if conf.Get(e.configPrefix+"trace.jaeger") != nil {
		var config = jaeger.RawConfig(e.configPrefix + "trace.jaeger")
		trace.SetGlobalTracer(config.Build())
	}
	return nil
}

// initMaxProcs init
func initMaxProcs() error {
	if maxProcs := conf.GetInt("maxProc"); maxProcs != 0 {
		runtime.GOMAXPROCS(maxProcs)
	} else {
		if _, err := maxprocs.Set(); err != nil {
			elog.EgoLogger.Panic("auto max procs", elog.FieldMod(ecode.ModProc), elog.FieldErrKind(ecode.ErrKindAny), elog.FieldErr(err))
		}
	}
	elog.EgoLogger.Info("auto max procs", elog.FieldMod(ecode.ModProc), elog.Int64("procs", int64(runtime.GOMAXPROCS(-1))))
	return nil
}

// printBanner init
func printBanner() error {
	const banner = `
 Welcome to Ego, starting application ...
`
	fmt.Println(xcolor.Green(banner))
	return nil
}

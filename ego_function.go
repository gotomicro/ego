package ego

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	sentinelmetrics "github.com/alibaba/sentinel-golang/metrics"
	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/econf/manager"
	"github.com/gotomicro/ego/core/eflag"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/esentinel"
	"github.com/gotomicro/ego/core/etrace"
	"github.com/gotomicro/ego/core/etrace/otel"
	"github.com/gotomicro/ego/core/util/xcolor"
	"github.com/gotomicro/ego/internal/retry"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/sync/errgroup"
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
		// 区分强制退出、优雅退出
		grace := s != syscall.SIGQUIT
		go func() {
			// todo 父节点传context待考虑
			stopCtx, cancel := context.WithTimeout(context.Background(), e.opts.stopTimeout)
			defer func() {
				signal.Stop(sig)
				cancel()
			}()
			_ = e.Stop(stopCtx, grace)
			<-stopCtx.Done()
			// 记录服务器关闭时候，由于关闭过慢，无法正常关闭，被强制cancel
			if errors.Is(stopCtx.Err(), context.DeadlineExceeded) {
				elog.Error("waitSignals stop context err", elog.FieldErr(stopCtx.Err()))
			}
		}()
		<-sig
		elog.Error("waitSignals quit")
		// 因为os.Signal长度为2，那么这里会阻塞住，如果发送两次信号量，强制退出
		os.Exit(128 + int(s.(syscall.Signal))) // second signal. Exit directly.
	}()
}

func (e *Ego) startServers(ctx context.Context) error {
	// start multi servers
	for _, s := range e.servers {
		s := s
		e.cycle.Run(func() (err error) {
			_ = s.Init()
			err = e.registerer.RegisterService(ctx, s.Info())
			if err != nil {
				e.logger.Error("register service err", elog.FieldComponent(s.PackageName()), elog.FieldComponentName(s.Name()), elog.FieldErr(err))
			}
			defer func() {
				_ = e.registerer.UnregisterService(ctx, s.Info())
			}()
			e.logger.Info("start server", elog.FieldComponent(s.PackageName()), elog.FieldComponentName(s.Name()), elog.FieldAddr(s.Info().Label()))
			defer e.logger.Info("stop server", elog.FieldComponent(s.PackageName()), elog.FieldComponentName(s.Name()), elog.FieldErr(err), elog.FieldAddr(s.Info().Label()))
			err = s.Start()
			return
		})
	}
	return nil
}

func (e *Ego) startOrderServers(ctx context.Context) (err error, isNeedStop bool) {
	// start order servers
	for _, s := range e.orderServers {
		s := s
		_ = s.Prepare()
		// 如果存在短时任务，那么只执行短时任务
		// 说明job在前面执行
		// 如果job执行完后，下面的操作需要stop
		if len(e.jobs) > 0 {
			return e.startJobs(), true
		}
		_ = s.Init()
		e.cycle.Run(func() (err error) {
			err = e.registerer.RegisterService(ctx, s.Info())
			if err != nil {
				e.logger.Error("register service err", elog.FieldComponent(s.PackageName()), elog.FieldComponentName(s.Name()), elog.FieldErr(err))
			}
			defer func() {
				_ = e.registerer.UnregisterService(ctx, s.Info())
			}()
			e.logger.Info("start order server", elog.FieldComponent(s.PackageName()), elog.FieldComponentName(s.Name()), elog.FieldAddr(s.Info().Label()))
			defer e.logger.Info("stop order server", elog.FieldComponent(s.PackageName()), elog.FieldComponentName(s.Name()), elog.FieldErr(err), elog.FieldAddr(s.Info().Label()))
			err = s.Start()
			return
		})
		isHealth := false
		for r := retry.Begin(); r.Continue(ctx); {
			// 检测server的health接口
			// 如果成功，那么就跳出循环
			if s.Health() {
				isHealth = true
				break
			}
		}
		if !isHealth {
			return fmt.Errorf("start order server fail,err:  " + s.Name()), true
		}

	}
	return nil, false
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
	// wrap jobs
	for _, runner := range e.jobs {
		runner := runner
		jobs = append(jobs, func() error {
			return runner.Start()
		})
	}

	eg := errgroup.Group{}
	for _, fn := range jobs {
		eg.Go(fn)
	}
	return eg.Wait()
}

// parseFlags init
func (e *Ego) parseFlags() error {
	if !e.opts.disableFlagConfig {
		eflag.Register(&eflag.StringFlag{
			Name:    "config",
			Usage:   "--config",
			EnvVar:  constant.EgoConfigPath,
			Default: constant.DefaultConfig,
			Action:  func(name string, fs *eflag.FlagSet) {},
		})
	}

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
		EnvVar:  constant.EnvAppHost,
		Default: "0.0.0.0",
		Action:  func(string, *eflag.FlagSet) {},
	})
	return eflag.ParseWithArgs(e.opts.arguments)
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
		e.opts.afterStopClean = append(e.opts.afterStopClean, elog.DefaultLogger.Flush)
	}

	if econf.Get(e.opts.configPrefix+"logger.ego") != nil {
		elog.EgoLogger = elog.Load(e.opts.configPrefix + "logger.ego").Build(elog.WithFileName(elog.EgoLoggerName))
		elog.EgoLogger.Info("reinit ego logger", elog.FieldComponent(elog.PackageName))
		e.opts.afterStopClean = append(e.opts.afterStopClean, elog.EgoLogger.Flush)
	}
	return nil
}

// initTracer init global tracer
func (e *Ego) initTracer() error {
	var (
		container *otel.Config
	)

	if econf.Get(e.opts.configPrefix+"trace") != nil {
		container = otel.Load(e.opts.configPrefix + "trace")
	} else {
		// 设置默认trace
		container = otel.DefaultConfig()
	}

	// 禁用trace
	if econf.GetBool(e.opts.configPrefix + "trace.disable") {
		elog.EgoLogger.Info("disable trace", elog.FieldComponent("app"))
		return nil
	}

	tracer := container.Build()
	etrace.SetGlobalTracer(tracer)
	e.opts.afterStopClean = append(e.opts.afterStopClean, container.Stop)
	elog.EgoLogger.Info("init trace", elog.FieldComponent("app"))
	return nil
}

// initSentinel 启动sentinel
func (e *Ego) initSentinel() error {
	if econf.Get(e.opts.configPrefix+"sentinel") != nil {
		esentinel.Load(e.opts.configPrefix + "sentinel").Build()
		sentinelmetrics.RegisterSentinelMetrics(prometheus.DefaultRegisterer.(*prometheus.Registry))
	}
	return nil
}

// initMaxProcs init
func initMaxProcs() error {
	if maxProcs := econf.GetInt("ego.maxProc"); maxProcs != 0 {
		runtime.GOMAXPROCS(maxProcs)
	} else {
		if _, err := maxprocs.Set(); err != nil {
			elog.EgoLogger.Error("init max procs", elog.FieldComponent("app"), elog.FieldErr(err))
		}
	}
	elog.EgoLogger.Info("init app", elog.FieldComponent("app"), elog.Int("pid", os.Getpid()), elog.Int("coreNum", runtime.GOMAXPROCS(-1)))
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

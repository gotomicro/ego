package ego

import (
	"flag"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/eflag"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/task/ejob"
)

func Test_loadConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := loadConfig(); (err != nil) != tt.wantErr {
				t.Errorf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_startJobsNoJob(t *testing.T) {
	app := &Ego{}
	err := app.startJobs()
	assert.NoError(t, err)
}

func Test_startJobsOneJobErrNil(t *testing.T) {
	app := &Ego{
		jobs:   make(map[string]ejob.Ejob),
		logger: elog.EgoLogger,
	}
	app.Job(ejob.Job("test", func(context ejob.Context) error {
		return nil
	}))

	err := app.startJobs()
	assert.NoError(t, err)
}

func Test_startJobsOneJobErrNotNil(t *testing.T) {
	resetFlagSet()
	eflag.Register(
		&eflag.StringFlag{
			Name:    "job",
			Usage:   "--job",
			Default: "",
		},
	)
	err := eflag.Parse()
	assert.NoError(t, err)
	err1 := flag.Set("job", "test")
	assert.NoError(t, err1)

	app := &Ego{
		jobs:   make(map[string]ejob.Ejob),
		logger: elog.EgoLogger,
	}
	app.Job(ejob.Job("test", func(context ejob.Context) error {
		return fmt.Errorf("test")
	}))

	err2 := app.startJobs()
	assert.Equal(t, "test", err2.Error())
}

func resetFlagSet() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagObj := eflag.NewFlagSet(flag.CommandLine)
	flag.Bool("test.v", false, "verbose: print additional output")
	flag.Bool("test.paniconexit0", false, "panic on call to os.Exit(0)")
	flag.String("test.run", "", "run only tests and examples matching `regexp`")
	flag.String("test.testlogfile", "", "write test action log to `file` (for use only by cmd/go)")
	flag.String("test.coverprofile", "", "write a coverage profile to `file`")
	flag.String("test.outputdir", "", "write profiles to `dir`")
	flag.Uint("test.count", 1, "run tests and benchmarks `n` times")
	flag.String("test.list", "", "list tests, examples, and benchmarks matching `regexp` then exit")
	flag.String("test.memprofile", "", "write an allocation profile to `file`")
	flag.Int("test.memprofilerate", 0, "set memory allocation profiling `rate` (see runtime.MemProfileRate)")
	flag.String("test.cpuprofile", "", "write a cpu profile to `file`")
	flag.String("test.blockprofile", "", "write a goroutine blocking profile to `file`")
	flag.Int("test.blockprofilerate", 1, "set blocking profile `rate` (see runtime.SetBlockProfileRate)")
	flag.String("test.mutexprofile", "", "write a mutex contention profile to the named file after execution")
	flag.Int("test.mutexprofilefraction", 1, "if >= 0, calls runtime.SetMutexProfileFraction()")
	flag.String("test.trace", "", "write an execution trace to `file`")
	flag.Duration("test.timeout", 0, "panic test binary after duration `d` (default 0, timeout disabled)")
	flag.String("test.cpu", "", "comma-separated `list` of cpu counts to run each test with")
	flag.String("test.gocoverdir", "", "gocoverdir dir")
	flag.Int("test.parallel", runtime.GOMAXPROCS(0), "run at most `n` tests in parallel")
	eflag.SetFlagSet(flagObj)
}

func Test_runSerialFuncReturnError(t *testing.T) {
	args := []func() error{func() error {
		return nil
	}}
	err := runSerialFuncReturnError(args)
	assert.NoError(t, err)

	args2 := []func() error{func() error {
		return fmt.Errorf("error")
	}}
	err2 := runSerialFuncReturnError(args2)
	assert.EqualError(t, err2, "error")
}

func Test_runSerialFuncLogError(t *testing.T) {
	args := []func() error{func() error {
		return fmt.Errorf("Test_runSerialFuncLogError")
	}}
	runSerialFuncLogError(args)
	err := elog.EgoLogger.Flush()
	assert.NoError(t, err)
	filePath := path.Join(elog.EgoLogger.ConfigDir(), elog.EgoLogger.ConfigName())
	logged, err1 := os.ReadFile(filePath)
	assert.NoError(t, err1)
	assert.Contains(t, string(logged), `"Test_runSerialFuncLogError"`)
}

func Test_initLogger(t *testing.T) {
	app := &Ego{}
	err := os.Setenv(constant.EgoDebug, "true")
	assert.NoError(t, err)
	cfg := `
[logger.default]
   debug = true
   enableAddCaller = true
`
	err1 := econf.LoadFromReader(strings.NewReader(cfg), toml.Unmarshal)
	assert.NoError(t, err1)

	err2 := app.initLogger()
	assert.NoError(t, err2)
	elog.Info("hello")
	err3 := elog.DefaultLogger.Flush()
	assert.NoError(t, err3)
	filePath := path.Join(elog.DefaultLogger.ConfigDir(), elog.DefaultLogger.ConfigName())
	logged, err4 := os.ReadFile(filePath)
	assert.NoError(t, err4)
	// 验证日志打印的caller是否正确 当前位置为ego/ego_function_test.go:150
	assert.Contains(t, string(logged), "hello", `ego/ego_function_test.go:150`)
}

func Test_initSysLogger(t *testing.T) {
	t.Run("没有ego的配置内容", func(t *testing.T) {
		app := &Ego{}
		cfg := ``
		err := econf.LoadFromReader(strings.NewReader(cfg), toml.Unmarshal)
		assert.NoError(t, err)

		err1 := app.initLogger()
		assert.NoError(t, err1)
		elog.EgoLogger.Info("hello1")
		err2 := elog.EgoLogger.Flush()
		assert.NoError(t, err2)
		filePath := path.Join(elog.EgoLogger.ConfigDir(), elog.EgoLogger.ConfigName())
		logged, err3 := os.ReadFile(filePath)
		assert.NoError(t, err3)
		// 验证日志是否打印了hello
		assert.Contains(t, string(logged), "hello1")
		// 验证日志文件名是否为ego.sys
		assert.Equal(t, elog.EgoLoggerName, elog.EgoLogger.ConfigName())
	})

	t.Run("有ego的配置内容，但是没有配置name选项", func(t *testing.T) {
		econf.Reset()
		app := &Ego{}
		cfg := `
[logger.ego]
   debug = true
`
		err := econf.LoadFromReader(strings.NewReader(cfg), toml.Unmarshal)
		assert.NoError(t, err)

		err1 := app.initLogger()
		assert.NoError(t, err1)
		elog.EgoLogger.Info("hello2")
		err2 := elog.EgoLogger.Flush()
		assert.NoError(t, err2)
		filePath := path.Join(elog.EgoLogger.ConfigDir(), elog.EgoLogger.ConfigName())
		logged, err3 := os.ReadFile(filePath)
		assert.NoError(t, err3)
		// 验证日志是否打印了hello
		assert.Contains(t, string(logged), "hello2")
		// 验证日志文件名是否为ego.sys
		assert.Equal(t, elog.EgoLoggerName, elog.EgoLogger.ConfigName())
	})

	t.Run("有ego的配置内容，并且有配置name选项", func(t *testing.T) {
		econf.Reset()
		app := &Ego{}
		fileName := "ego.sys.log"
		cfg := `
[logger.ego]
   debug = true
   name = "ego.sys.log"
`
		err := econf.LoadFromReader(strings.NewReader(cfg), toml.Unmarshal)
		assert.NoError(t, err)

		err1 := app.initLogger()
		assert.NoError(t, err1)
		elog.EgoLogger.Info("hello3")
		err2 := elog.EgoLogger.Flush()
		assert.NoError(t, err2)
		filePath := path.Join(elog.EgoLogger.ConfigDir(), elog.EgoLogger.ConfigName())
		logged, err3 := os.ReadFile(filePath)
		assert.NoError(t, err3)
		// 验证日志是否打印了hello
		assert.Contains(t, string(logged), "hello3")
		// 验证日志文件名是否为ego.sys.log
		assert.Equal(t, fileName, elog.EgoLogger.ConfigName())
	})

	t.Run("修改EgoLogger本身，ego中的logger同步生效", func(t *testing.T) {
		econf.Reset()
		// 先还原一下默认的EgoLogger
		elog.EgoLogger = elog.DefaultContainer().Build(elog.WithFileName(elog.EgoLoggerName))
		var (
			app = &Ego{
				logger: elog.EgoLogger, // logger与ego.New()方法中保持一致
			}
			cfg = `
				[logger.ego]
				   debug = true
				   name = "ego.sys.log" 
				`
		)

		err := econf.LoadFromReader(strings.NewReader(cfg), toml.Unmarshal)
		assert.NoError(t, err)

		// 初始化之前使用的默认的name
		assert.Equal(t, elog.EgoLoggerName, app.logger.ConfigName())

		err1 := app.initLogger()
		assert.NoError(t, err1)

		// 初始化后验证ego结构体中的logger日志文件名是否为ego.sys.log
		assert.Equal(t, "ego.sys.log", app.logger.ConfigName())
	})
}

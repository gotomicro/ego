package ego

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/eflag"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/task/ejob"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
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
	assert.Nil(t, err)
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
	assert.Nil(t, err)
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
	err = flag.Set("job", "test")
	assert.NoError(t, err)

	app := &Ego{
		jobs:   make(map[string]ejob.Ejob),
		logger: elog.EgoLogger,
	}
	app.Job(ejob.Job("test", func(context ejob.Context) error {
		return fmt.Errorf("test")
	}))

	err = app.startJobs()
	assert.Equal(t, "test", err.Error())
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
	flag.Int("test.parallel", runtime.GOMAXPROCS(0), "run at most `n` tests in parallel")
	eflag.SetFlagSet(flagObj)
}

func Test_runSerialFuncReturnError(t *testing.T) {
	args := []func() error{func() error {
		return nil
	}}
	err := runSerialFuncReturnError(args)
	assert.Nil(t, err)

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
	elog.EgoLogger.Flush()
	filePath := path.Join(elog.EgoLogger.ConfigDir(), elog.EgoLogger.ConfigName())
	logged, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Contains(t, string(logged), `"Test_runSerialFuncLogError"`)
}

func Test_initLogger(t *testing.T) {
	app := &Ego{}
	err := os.Setenv(constant.EgoDebug, "true")
	assert.Nil(t, err)
	cfg := `
[logger.default]
   debug = true
   enableAddCaller = true
`
	err = econf.LoadFromReader(strings.NewReader(cfg), toml.Unmarshal)
	assert.NoError(t, err)

	err = app.initLogger()
	assert.Nil(t, err)
	elog.Info("hello")
	elog.DefaultLogger.Flush()
	filePath := path.Join(elog.DefaultLogger.ConfigDir(), elog.DefaultLogger.ConfigName())
	logged, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	// 验证日志打印的caller是否正确 当前位置为ego/ego_function_test.go:150
	assert.Contains(t, string(logged), "hello", `ego/ego_function_test.go:150`)
}

func Test_initSysLogger(t *testing.T) {
	t.Run("没有ego的配置内容", func(t *testing.T) {
		app := &Ego{}
		cfg := ``
		err := econf.LoadFromReader(strings.NewReader(cfg), toml.Unmarshal)
		assert.NoError(t, err)

		err = app.initLogger()
		assert.Nil(t, err)
		elog.EgoLogger.Info("hello1")
		elog.EgoLogger.Flush()
		filePath := path.Join(elog.EgoLogger.ConfigDir(), elog.EgoLogger.ConfigName())
		logged, err := os.ReadFile(filePath)
		assert.Nil(t, err)
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

		err = app.initLogger()
		assert.Nil(t, err)
		elog.EgoLogger.Info("hello2")
		elog.EgoLogger.Flush()
		filePath := path.Join(elog.EgoLogger.ConfigDir(), elog.EgoLogger.ConfigName())
		logged, err := os.ReadFile(filePath)
		assert.Nil(t, err)
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

		err = app.initLogger()
		assert.Nil(t, err)
		elog.EgoLogger.Info("hello3")
		elog.EgoLogger.Flush()
		filePath := path.Join(elog.EgoLogger.ConfigDir(), elog.EgoLogger.ConfigName())
		logged, err := os.ReadFile(filePath)
		assert.Nil(t, err)
		// 验证日志是否打印了hello
		assert.Contains(t, string(logged), "hello3")
		// 验证日志文件名是否为ego.sys.log
		assert.Equal(t, fileName, elog.EgoLogger.ConfigName())
	})
}

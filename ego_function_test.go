package ego

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/gotomicro/ego/core/eflag"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/task/ejob"
	"github.com/stretchr/testify/assert"
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

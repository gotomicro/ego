package eflag

import (
	"flag"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/constant"
)

const (
	ConfigEnvToml     = "config/env.toml"
	ConfigDefaultToml = "config/default.toml"
	ConfigFlagToml    = "config/flag.toml"
)

func TestFlagSet_Register_Length(t *testing.T) {
	resetFlagSet()
	Register(&StringFlag{
		Name:    "config",
		Usage:   "--config",
		EnvVar:  constant.EgoConfigPath,
		Default: ConfigDefaultToml,
		Action:  func(name string, fs *FlagSet) {},
	})
	assert.Equal(t, 1, len(flagset.flags))
}

func TestFlagSetInt_Register_Length(t *testing.T) {
	resetFlagSet()
	Register(&IntFlag{
		Name:   "int",
		Usage:  "--int",
		Action: func(name string, fs *FlagSet) {},
	})
	assert.Equal(t, 1, len(flagset.flags))
}

func TestFlagSetUint_Register_Length(t *testing.T) {
	resetFlagSet()
	Register(&UintFlag{
		Name:   "uint",
		Usage:  "--uint",
		Action: func(name string, fs *FlagSet) {},
	})
	assert.Equal(t, 1, len(flagset.flags))
}

func TestFlagSetFloat64_Register_Length(t *testing.T) {
	resetFlagSet()
	Register(&UintFlag{
		Name:   "float64",
		Usage:  "--float64",
		Action: func(name string, fs *FlagSet) {},
	})
	assert.Equal(t, 1, len(flagset.flags))
}
func TestFlagSet_Register_Default(t *testing.T) {
	resetFlagSet()
	Register(&StringFlag{
		Name:    "config",
		Usage:   "--config",
		EnvVar:  constant.EgoConfigPath,
		Default: ConfigDefaultToml,
		Action:  func(name string, fs *FlagSet) {},
	})
	err := Parse()
	assert.NoError(t, err)

	configStr, err := StringE("config")
	assert.NoError(t, err)
	assert.Equal(t, ConfigDefaultToml, configStr)
}

func TestFlagSet_Register_Env(t *testing.T) {
	os.Setenv(constant.EgoConfigPath, "config/env.toml")
	defer os.Unsetenv(constant.EgoConfigPath)
	resetFlagSet()

	Register(&StringFlag{
		Name:   "config",
		Usage:  "--config",
		EnvVar: constant.EgoConfigPath,
		Action: func(name string, fs *FlagSet) {},
	})
	err := Parse()
	assert.NoError(t, err)

	configStr, err := StringE("config")
	assert.NoError(t, err)
	assert.Equal(t, ConfigEnvToml, configStr)
}

func TestFlagSet_Register_Flag(t *testing.T) {
	os.Setenv(constant.EgoConfigPath, "config/env.toml")
	defer os.Unsetenv(constant.EgoConfigPath)
	resetFlagSet()

	Register(&StringFlag{
		Name:   "config",
		Usage:  "--config",
		EnvVar: constant.EgoConfigPath,
		Action: func(name string, fs *FlagSet) {},
	})
	err := Parse()
	assert.NoError(t, err)

	err = flag.Set("config", ConfigFlagToml)
	assert.NoError(t, err)
	configStr, err := StringE("config")
	assert.NoError(t, err)
	assert.Equal(t, ConfigFlagToml, configStr)
}

func TestFlagSet_Register_Priority(t *testing.T) {
	// 1 设置了 flag，env，default config，那么应该为flag config
	_ = os.Setenv(constant.EgoConfigPath, "config/env.toml")
	defer os.Unsetenv(constant.EgoConfigPath)
	resetFlagSet()
	Register(&StringFlag{
		Name:    "config",
		Usage:   "--config",
		EnvVar:  constant.EgoConfigPath,
		Default: ConfigDefaultToml,
		Action:  func(name string, fs *FlagSet) {},
	})
	_ = Parse()
	_ = flag.Set("config", ConfigFlagToml)
	configStr, err := StringE("config")
	assert.NoError(t, err)
	assert.Equal(t, ConfigFlagToml, configStr)

	// 2 设置了 env，default config，那么应该为env config
	_ = os.Setenv(constant.EgoConfigPath, "config/env.toml")
	resetFlagSet()
	Register(&StringFlag{
		Name:    "config",
		Usage:   "--config",
		EnvVar:  constant.EgoConfigPath,
		Default: ConfigDefaultToml,
		Action:  func(name string, fs *FlagSet) {},
	})
	_ = Parse()
	configStr, err = StringE("config")
	assert.NoError(t, err)
	assert.Equal(t, ConfigEnvToml, configStr)

	// 3 设置了 default config，那么应该为default config
	os.Unsetenv(constant.EgoConfigPath)
	resetFlagSet()
	Register(&StringFlag{
		Name:    "config",
		Usage:   "--config",
		EnvVar:  constant.EgoConfigPath,
		Default: ConfigDefaultToml,
		Action:  func(name string, fs *FlagSet) {},
	})
	err = Parse()
	assert.NoError(t, err)
	configStr, err = StringE("config")
	assert.NoError(t, err)
	assert.Equal(t, ConfigDefaultToml, configStr)
}

func TestFlagSet_Register_Bool(t *testing.T) {
	resetFlagSet()
	Register(&BoolFlag{
		Name:    "watch",
		Usage:   "--watch",
		Default: true,
		Action:  func(name string, fs *FlagSet) {},
	})
	err := Parse()
	assert.NoError(t, err)
	boolFlag, err := BoolE("watch")
	assert.NoError(t, err)
	assert.Equal(t, true, boolFlag)

	os.Setenv("EGO_WATCH", "false")
	defer os.Unsetenv("EGO_WATCH")
	resetFlagSet()
	Register(&BoolFlag{
		Name:    "watch",
		Usage:   "--watch",
		Default: true,
		EnvVar:  "EGO_WATCH",
		Action:  func(name string, fs *FlagSet) {},
	})
	err = Parse()
	assert.NoError(t, err)
	boolFlag, err = BoolE("watch")
	assert.NoError(t, err)
	assert.Equal(t, false, boolFlag)

	resetFlagSet()
	Register(&BoolFlag{
		Name:    "watch",
		Usage:   "--watch",
		Default: true,
		Action:  func(name string, fs *FlagSet) {},
	})
	err = Parse()
	assert.NoError(t, err)
	err = flag.Set("watch", "true")
	assert.NoError(t, err)
	boolFlag, err = BoolE("watch")
	assert.NoError(t, err)
	assert.Equal(t, true, boolFlag)

	resetFlagSet()
	Register(&BoolFlag{
		Name:    "watch",
		Usage:   "--watch",
		Default: true,
		Action:  func(name string, fs *FlagSet) {},
	})
	err = Parse()
	assert.NoError(t, err)
	err = flag.Set("watch", "false")
	assert.NoError(t, err)
	boolFlag, err = BoolE("watch")
	assert.NoError(t, err)
	assert.Equal(t, false, boolFlag)
}

func TestNewFlagSet(t *testing.T) {
	obj := NewFlagSet(nil)
	assert.Nil(t, obj.FlagSet)
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	obj2 := NewFlagSet(flag.CommandLine)
	assert.True(t, assert.ObjectsAreEqual(flag.CommandLine, obj2.FlagSet))
}

func TestParseWithArgs(t *testing.T) {
	resetFlagSet()
	Register(&BoolFlag{
		Name:    "bool",
		Usage:   "--bool",
		Default: false,
		Action:  func(name string, fs *FlagSet) {},
	})
	err := ParseWithArgs([]string{"--bool"})
	assert.NoError(t, err)
	boolFlag, err := BoolE("bool")
	assert.NoError(t, err)
	assert.Equal(t, true, boolFlag)

	resetFlagSet()
	Register(&BoolFlag{
		Name:    "bool",
		Usage:   "--bool",
		Default: true,
		Action:  func(name string, fs *FlagSet) {},
	})
	err = ParseWithArgs([]string{"--bool=false"})
	assert.NoError(t, err)
	boolFlag, err = BoolE("bool")
	assert.NoError(t, err)
	assert.Equal(t, false, boolFlag)

	resetFlagSet()
	Register(&StringFlag{
		Name:    "string",
		Usage:   "--string",
		Default: "world",
		Action:  func(name string, fs *FlagSet) {},
	})
	err = ParseWithArgs([]string{"--string", "hello"})
	assert.NoError(t, err)
	strFlag, err := StringE("string")
	assert.NoError(t, err)
	assert.Equal(t, "hello", strFlag)
}

func resetFlagSet() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagObj := NewFlagSet(flag.CommandLine)
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
	flag.String("test.gocoverdir", "", "gocoverdir dir")
	SetFlagSet(flagObj)
}

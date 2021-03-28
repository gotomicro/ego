package eflag

import (
	"flag"
	"os"
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

func TestFlagSet_Register_Default(t *testing.T) {
	resetFlagSet()
	Register(&StringFlag{
		Name:    "config",
		Usage:   "--config",
		EnvVar:  constant.EgoConfigPath,
		Default: ConfigDefaultToml,
		Action:  func(name string, fs *FlagSet) {},
	})
	Parse()
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
	Parse()
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
	Parse()

	flag.Set("config", ConfigFlagToml)
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
	Parse()
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
	Parse()
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
	Parse()
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
	Parse()
	flag.Set("watch", "true")
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
	Parse()
	flag.Set("watch", "false")
	boolFlag, err = BoolE("watch")
	assert.NoError(t, err)
	assert.Equal(t, false, boolFlag)
}

func resetFlagSet() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagObj := &FlagSet{
		FlagSet: flag.CommandLine,
		actions: make(map[string]func(string, *FlagSet)),
	}
	flag.Bool("test.v", false, "verbose: print additional output")
	flag.Bool("test.paniconexit0", false, "panic on call to os.Exit(0)")
	flag.String("test.run", "", "run only tests and examples matching `regexp`")
	flag.Duration("test.timeout", 0, "panic test binary after duration `d` (default 0, timeout disabled)")
	flag.String("test.testlogfile", "", "write test action log to `file` (for use only by cmd/go)")
	setFlagSet(flagObj)
}

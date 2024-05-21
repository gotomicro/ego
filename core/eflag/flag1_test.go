package eflag

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/constant"
)

func TestApply(t *testing.T) {
	_ = os.Setenv(constant.EgoConfigPath, "config/env.toml")
	defer os.Unsetenv(constant.EgoConfigPath)
	resetFlagSet()

	Register(&Float64Flag{
		Name:    "watch",
		Usage:   "--watch",
		Default: 222,
		EnvVar:  constant.EgoConfigPath,
		Action:  func(name string, fs *FlagSet) {},
	})
	_ = Parse()
	_ = flag.Set("config", ConfigFlagToml)

	err := ParseWithArgs([]string{"--watch-false"})
	assert.NoError(t, err)

	Float64("watch")
	assert.NoError(t, nil)
	_, err1 := Float64E("watch")
	assert.NoError(t, err1)

}

func TestInt(t *testing.T) {
	_ = os.Setenv(constant.EgoConfigPath, "config/env.toml")
	defer os.Unsetenv(constant.EgoConfigPath)
	resetFlagSet()
	Register(&IntFlag{
		Name:    "watch",
		Usage:   "--watch",
		Default: 222,
		EnvVar:  constant.EgoConfigPath,
		Action:  func(name string, fs *FlagSet) {},
	})
	_ = Parse()
	_ = flag.Set("config", ConfigFlagToml)
	err := ParseWithArgs([]string{"--watch-false"})
	assert.NoError(t, err)
	Int("watch")
	assert.NoError(t, nil)
	_, err1 := IntE("watch")
	assert.NoError(t, err1)
}

func TestUint(t *testing.T) {
	_ = os.Setenv(constant.EgoConfigPath, "config/env.toml")
	defer os.Unsetenv(constant.EgoConfigPath)
	resetFlagSet()
	Register(&UintFlag{
		Name:    "watch",
		Usage:   "--watch",
		Default: 222,
		EnvVar:  constant.EgoConfigPath,
		Action:  func(name string, fs *FlagSet) {},
	})
	_ = Parse()
	_ = flag.Set("config", ConfigFlagToml)
	err := ParseWithArgs([]string{"--watch-false"})
	assert.NoError(t, err)
	Uint("watch")
	assert.NoError(t, nil)
	_, err1 := UintE("watch")
	assert.NoError(t, err1)
}

func TestString(t *testing.T) {
	_ = os.Setenv(constant.EgoConfigPath, "config/env.toml")
	defer os.Unsetenv(constant.EgoConfigPath)
	resetFlagSet()
	Register(&StringFlag{
		Name:    "watch",
		Usage:   "--watch",
		Default: "test",
		EnvVar:  constant.EgoConfigPath,
		Action:  func(name string, fs *FlagSet) {},
	})
	_ = Parse()
	_ = flag.Set("config", ConfigFlagToml)
	err := ParseWithArgs([]string{"--watch-false"})
	assert.NoError(t, err)
	String("watch")
	assert.NoError(t, nil)
	_, err1 := StringE("watch")
	assert.NoError(t, err1)
}

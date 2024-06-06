package eflag

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/constant"
)

func TestApply(t *testing.T) {
	err1 := os.Setenv(constant.EgoConfigPath, "config/env.toml")
	assert.NoError(t, err1)
	defer os.Unsetenv(constant.EgoConfigPath)
	resetFlagSet()

	Register(&Float64Flag{
		Name:    "watch",
		Usage:   "--watch",
		Default: 222,
		EnvVar:  constant.EgoConfigPath,
		Action:  func(name string, fs *FlagSet) {},
	})
	err2 := Parse()
	assert.NoError(t, err2)
	_ = flag.Set("config", ConfigFlagToml)

	err := ParseWithArgs([]string{"--watch-false"})
	assert.NoError(t, err)
	assert.Equal(t, float64(0), Float64("watch"))
	out, err3 := Float64E("watch")
	assert.NoError(t, err3)
	assert.Equal(t, float64(0), out)
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
	assert.Equal(t, int64(0), Int("watch"))
	_, err1 := IntE("watch")
	assert.NoError(t, err1)
}

func TestUint(t *testing.T) {
	err1 := os.Setenv(constant.EgoConfigPath, "config/env.toml")
	assert.NoError(t, err1)
	defer os.Unsetenv(constant.EgoConfigPath)
	resetFlagSet()
	Register(&UintFlag{
		Name:    "watch",
		Usage:   "--watch",
		Default: 222,
		EnvVar:  constant.EgoConfigPath,
		Action:  func(name string, fs *FlagSet) {},
	})
	err2 := Parse()
	assert.NoError(t, err2)
	_ = flag.Set("config", ConfigFlagToml)
	err := ParseWithArgs([]string{"--watch-false"})
	assert.NoError(t, err)
	Uint("watch")
	assert.Equal(t, uint64(0), Uint("watch"))
	out, err1 := UintE("watch")
	assert.NoError(t, err1)
	assert.Equal(t, uint64(0), out)
}

func TestString(t *testing.T) {
	err := os.Setenv(constant.EgoConfigPath, "config/env.toml")
	assert.NoError(t, err)
	defer os.Unsetenv(constant.EgoConfigPath)
	resetFlagSet()
	Register(&StringFlag{
		Name:    "watch",
		Usage:   "--watch",
		Default: "test",
		EnvVar:  constant.EgoConfigPath,
		Action:  func(name string, fs *FlagSet) {},
	})
	err1 := Parse()
	assert.NoError(t, err1)
	_ = flag.Set("config", ConfigFlagToml)
	err2 := ParseWithArgs([]string{"--watch-false"})
	assert.NoError(t, err2)
	assert.Equal(t, "config/env.toml", String("watch"))
	out, err3 := StringE("watch")
	assert.NoError(t, err3)
	assert.Equal(t, "config/env.toml", out)
}

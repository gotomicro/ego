package elog

import (
	"os"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/econf"
)

func TestRotateLogger(t *testing.T) {
	err := os.Setenv("EGO_DEBUG", "false")
	assert.NoError(t, err)
	conf := `
[default]
debug = false
level = "info"
enableAsync = false
`
	err = econf.LoadFromReader(strings.NewReader(conf), toml.Unmarshal)
	assert.NoError(t, err)
	cmp := Load("default").Build().With(String("prefix", "PREFIX"))
	defer cmp.Flush()
	cmp.Error("TestRotateLogger test", String("name", "lee"), Int("age", 17))
}

func TestAlislsLogger(t *testing.T) {
	err := os.Setenv("EGO_DEBUG", "false")
	assert.NoError(t, err)
	conf := `
[sls]
level = "info"
enableAsync = false
writer = "ali"
`
	err = econf.LoadFromReader(strings.NewReader(conf), toml.Unmarshal)
	assert.NoError(t, err)
	cmp := Load("sls").Build()
	newCmp := cmp.With(String("prefix", "PREFIX"))
	defer newCmp.Flush()
	cmp.Error("TestRotateLogger test", String("name", "lee"), Int("age", 17))
	cmp.Error("TestRotateLogger test", String("name", "lee"), Int("age", 17))
	cmp.Error("TestRotateLogger test", String("name", "lee"), Int("age", 17))
	// time.Sleep(6 * time.Second)
}

package ehttp

import (
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
)

func TestLoad(t *testing.T) {
	file, err := os.Open("./config_test/conf.toml")
	assert.NoError(t, err)
	err1 := econf.LoadFromReader(file, toml.Unmarshal)
	assert.NoError(t, err1)
	Load("test").Build()
	logger := DefaultContainer().logger.With(elog.FieldComponentName("test"))
	assert.Equal(t, logger, Load("test").logger)
}

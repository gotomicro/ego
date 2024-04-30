package ehttp

import (
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/econf"
)

func TestLoad(t *testing.T) {
	file, err := os.Open("./config_test/conf.toml")
	assert.NoError(t, err)
	err = econf.LoadFromReader(file, toml.Unmarshal)
	assert.NoError(t, err)
	container := Load("test").Build().name
	assert.Equal(t, "test", container)
}

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
	logger := Load("default").Build()
	defer logger.Flush()
	logger.Error("1")

	childLogger := logger.With(String("prefix", "PREFIX"))
	childLogger.Error("2")
	defer childLogger.Flush()

	logger.Error("3")
	logger.With(String("prefix2", "PREFIX2"))
	logger.Error("4")
}

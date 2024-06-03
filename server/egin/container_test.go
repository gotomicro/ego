package egin

import (
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
)

func TestLoadAndBuild(t *testing.T) {
	conf := `[test]
AccessInterceptorReqResFilter = "test"
EnableLocalMainIP = true
EnableTraceInterceptor = true
EnableSentinel = true`
	err := econf.LoadFromReader(strings.NewReader(conf), toml.Unmarshal)
	assert.NoError(t, err)
	load := Load("test")
	logger := DefaultContainer().logger.With(elog.FieldComponentName("test"))
	assert.Equal(t, logger, load.logger)
	c := load.Build()
	assert.NotNil(t, c)
}

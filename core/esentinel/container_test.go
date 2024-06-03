package esentinel

import (
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
)

func TestDefaultContainer(t *testing.T) {
	in := &Container{
		config: DefaultConfig(),
		logger: elog.EgoLogger.With(elog.FieldComponent(PackageName)),
	}
	out := DefaultContainer()
	assert.Equal(t, in, out)
}

func TestLoad(t *testing.T) {
	conf := `
[test]
addr = ":9091"
`
	err := econf.LoadFromReader(strings.NewReader(conf), toml.Unmarshal)
	assert.NoError(t, err)
	Load("test").Build()
	in := &Container{
		name:   "test",
		config: DefaultConfig(),
		logger: DefaultContainer().logger.With(elog.FieldComponentName("test")),
	}
	assert.Equal(t, in, Load("test"))
}

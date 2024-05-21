package egovernor

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
	assert.NoError(t, nil)
}

func TestBuild(t *testing.T) {
	var c = &Container{
		name:   "test",
		config: &Config{Host: "test", Port: 8080},
		logger: nil,
	}
	opt1 := WithHost("test")
	opt2 := WithPort(8080)
	c.Build(opt1, opt2)
	assert.NoError(t, nil)
}

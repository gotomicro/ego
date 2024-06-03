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

func TestLoadAndBuild(t *testing.T) {
	conf := `
[test]
Host = "172.16.21.157"
Port = 8080
EnableLocalMainIP = true
Network = "tcp4"
`
	err := econf.LoadFromReader(strings.NewReader(conf), toml.Unmarshal)
	assert.NoError(t, err)
	l := Load("test")
	c := &Container{
		config: &Config{
			Host:              "172.16.21.157",
			Port:              8080,
			EnableLocalMainIP: true,
			Network:           "tcp4",
		},
		name:   "test",
		logger: DefaultContainer().logger.With(elog.FieldComponentName("test")),
	}
	// assert.Equal(t, c, l)
	assert.Equal(t, c.name, l.name)
	opt1 := WithHost("172.16.21.157")
	opt2 := WithPort(8080)
	e := l.Build(opt1, opt2)
	assert.Equal(t, c.name, e.name)
	assert.Equal(t, c.logger, e.logger)
	assert.Equal(t, c.config, e.config)
}

package ecron

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/robfig/cron/v3"

	"github.com/gotomicro/ego/core/econf"
)

func container() *Container {
	err := econf.LoadFromReader(strings.NewReader(`[cron.test]
spec = "0 0 1 1 *"`), toml.Unmarshal)
	if err != nil {
		panic(err.Error())
	}
	return Load("cron.test")
}

func TestWithLock(t *testing.T) {
	lock := &mockLock{}
	comp := container().Build(WithLock(lock))
	if comp.config.lock != lock {
		t.Failed()
	}
}

func TestWithWrappers(t *testing.T) {
	a := 0
	wrapper := func(job cron.Job) cron.Job {
		a = 1
		return job
	}
	comp := container().Build(WithWrappers(wrapper))
	wrapperLen := len(comp.config.wrappers)
	comp.config.wrappers[wrapperLen-1](nil)
	if a != 1 {
		t.Failed()
	}
}

func TestWithSeconds(t *testing.T) {
	c := &Container{
		config: &Config{
			EnableSeconds: false,
		},
	}
	WithSeconds()(c)
	if !c.config.EnableSeconds {
		t.Error("expect EnableSeconds = true")
	}
}

func TestWithParser(t *testing.T) {
	p := cron.NewParser(cron.Second | cron.Month)
	c := &Container{
		config: &Config{
			parser: cron.NewParser(cron.Second),
		},
	}
	WithParser(p)(c)
	if !reflect.DeepEqual(c.config.parser, p) {
		t.Error("expect parser equal")
	}
}

func TestWithLocation(t *testing.T) {
	c := &Container{
		config: DefaultConfig(),
	}
	loc := time.UTC
	WithLocation(loc)(c)
	if !reflect.DeepEqual(loc, c.config.loc) {
		t.Error("expect loc equal")
	}
}

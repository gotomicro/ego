package egovernor

import (
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/emetric"
	"github.com/gotomicro/ego/core/util/xnet"
)

// Container defines a component instance.
type Container struct {
	config *Config
	name   string
	err    error
	logger *elog.Component
}

// DefaultContainer returns an default container.
func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
		logger: elog.EgoLogger.With(elog.FieldComponent(PackageName)),
	}
}

// Load parses container configuration from configuration provider, such as a toml file,
// then use the configuration to construct a component container.
func Load(key string) *Container {
	c := DefaultContainer()
	c.logger = c.logger.With(elog.FieldComponentName(key))
	if err := econf.UnmarshalKey(key, &c.config); err != nil {
		c.err = err
		return c
	}
	var (
		host string
		err  error
	)
	// 获取网卡ip
	if c.config.EnableLocalMainIP {
		host, _, err = xnet.GetLocalMainIP()
		if err != nil {
			elog.Error("get local main ip error", elog.FieldErr(err))
		} else {
			c.config.Host = host
		}
	}
	c.name = key
	return c
}

// Build constructs a specific component from container.
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	if c.config.EnableConnTcpMetric {
		obj := emetric.NewTCPStatCollector(c.config.ConnTcpMetricPorts)
		obj.Update()
	}
	return newComponent(c.name, c.config, c.logger)
}

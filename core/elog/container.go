package elog

import (
	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/econf"
)

// Container defines a component instance.
type Container struct {
	config *Config
	name   string
}

// DefaultContainer returns an default container.
func DefaultContainer() *Container {
	return &Container{
		config: defaultConfig(),
	}
}

// Load parses container configuration from configuration provider, such as a toml file,
// then use the configuration to construct a component container.
func Load(key string) *Container {
	c := DefaultContainer()
	if err := econf.UnmarshalKey(key, &c.config); err != nil {
		panic(err)
	}
	c.name = key
	return c
}

// Build constructs a specific component from container.
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}

	if eapp.IsDevelopmentMode() {
		c.config.Debug = true           // 调试模式，终端输出
		c.config.EnableAsync = false    // 调试模式，同步输出
		c.config.EnableAddCaller = true // 调试模式，增加行号输出
	}

	// 如果用户设置了该配置，那么该配置被用户接管
	// 如果用户没有设置，那么使用默认配置
	if c.config.encoderConfig == nil {
		if c.config.Debug {
			c.config.encoderConfig = defaultDebugConfig()
		} else {
			c.config.encoderConfig = defaultZapConfig()
		}
	}

	if eapp.EnableLoggerAddApp() {
		c.config.fields = append(c.config.fields, FieldApp(eapp.Name()))
	}

	// 设置ego日志的log name，用于stderr区分系统日志和业务日志
	// config writer setting > env writer setting
	if c.config.Writer == "stderr" || (c.config.Writer == "" && eapp.EgoLogWriter() == "stderr") {
		c.config.fields = append(c.config.fields, FieldLogName(c.config.Name))
	}

	return newLogger(c.name, c.name, c.config)
}

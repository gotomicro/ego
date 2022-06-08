package elog

import (
	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/econf"
)

// Container 容器
type Container struct {
	config *Config
	name   string
}

// DefaultContainer 默认容器
func DefaultContainer() *Container {
	return &Container{
		config: defaultConfig(),
	}
}

// Load 加载配置key
func Load(key string) *Container {
	c := DefaultContainer()
	if err := econf.UnmarshalKey(key, &c.config); err != nil {
		panic(err)
	}
	c.name = key
	return c
}

// Build 构建组件
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}

	if eapp.IsDevelopmentMode() {
		c.config.Debug = true           // 调试模式，终端输出
		c.config.EnableAsync = false    // 调试模式，同步输出
		c.config.EnableAddCaller = true // 调试模式，增加行号输出
	}

	if c.config.Debug {
		c.config.encoderConfig = defaultDebugConfig()
	}

	if eapp.EnableLoggerAddApp() {
		c.config.fields = append(c.config.fields, FieldApp(eapp.Name()))
	}

	// 设置ego日志的log name，用于stderr区分系统日志和业务日志
	if eapp.EgoLogWriter() == "stderr" {
		c.config.fields = append(c.config.fields, FieldLogName(c.config.Name))
	}

	return newLogger(c.name, c.name, c.config)
}

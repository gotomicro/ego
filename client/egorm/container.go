package egorm

import (
	"github.com/gotomicro/ego/core/conf"
	"github.com/gotomicro/ego/core/ecode"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/metric"
)

type Option func(c *Container)

type Container struct {
	config *Config
	name   string
	logger *elog.Component
}

func DefaultContainer() *Container {
	return &Container{
		config: DefaultConfig(),
		logger: elog.EgoLogger.With(elog.FieldMod("client.egorm")),
	}
}

func Load(key string) *Container {
	c := DefaultContainer()
	if err := conf.UnmarshalKey(key, &c.config); err != nil {
		c.logger.Panic("parse config error", elog.FieldErr(err), elog.FieldKey(key))
		return c
	}

	c.name = key
	return c
}

// WithInterceptor ...
func WithInterceptor(is ...Interceptor) Option {
	return func(c *Container) {
		if c.config.interceptors == nil {
			c.config.interceptors = make([]Interceptor, 0)
		}
		c.config.interceptors = append(c.config.interceptors, is...)
	}
}

// Build ...
func (c *Container) Build(options ...Option) *Component {
	if options == nil {
		options = make([]Option, 0)
	}

	if c.config.Debug {
		options = append(options, WithInterceptor(debugInterceptor))
	}

	if !c.config.DisableTrace {
		options = append(options, WithInterceptor(traceInterceptor))
	}

	if !c.config.DisableMetric {
		options = append(options, WithInterceptor(metricInterceptor))
	}

	for _, option := range options {
		option(c)
	}

	var err error
	// todo 设置补齐超时时间
	// timeout 1s
	// readTimeout 5s
	// writeTimeout 5s
	c.config.dsnCfg, err = ParseDSN(c.config.DSN)

	if err == nil {
		c.logger.Info(ecode.MsgClientMysqlOpenStart, elog.FieldAddr(c.config.dsnCfg.Addr), elog.FieldName(c.config.dsnCfg.DBName))
	} else {
		c.logger.Panic(ecode.MsgClientMysqlOpenStart, elog.FieldErr(err))
	}

	component, err := newComponent(c.name, c.config, c.logger)
	if err != nil {
		if c.config.OnDialError == "panic" {
			c.logger.Panic("open mysql", elog.FieldErrKind(ecode.ErrKindRequestErr), elog.FieldErr(err), elog.FieldAddr(c.config.dsnCfg.Addr), elog.FieldValueAny(c.config))
		} else {
			metric.LibHandleCounter.Inc(metric.TypeGorm, c.name+".ping", c.config.dsnCfg.Addr, "open err")
			c.logger.Error("open mysql", elog.FieldErrKind(ecode.ErrKindRequestErr), elog.FieldErr(err), elog.FieldAddr(c.config.dsnCfg.Addr), elog.FieldValueAny(c.config))
			return component
		}
	}

	if err := component.DB.DB().Ping(); err != nil {
		c.logger.Panic("ping mysql", elog.FieldErrKind(ecode.ErrKindRequestErr), elog.FieldErr(err), elog.FieldValueAny(c.config))
	}

	// store db
	instances.Store(c.name, component)
	return component
}

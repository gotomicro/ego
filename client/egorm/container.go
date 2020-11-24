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
		logger: elog.EgoLogger.With(elog.FieldMod("client.egorm")),
	}
}

func Load(key string) *Container {
	c := DefaultContainer()
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		c.logger.Panic("parse config error", elog.FieldErr(err), elog.FieldKey(key))
		return c
	}

	c.config = config
	c.name = key
	return c
}

// WithInterceptor ...
func (c *Container) WithInterceptor(is ...Interceptor) *Container {
	if c.config.interceptors == nil {
		c.config.interceptors = make([]Interceptor, 0)
	}
	c.config.interceptors = append(c.config.interceptors, is...)
	return c
}

// Build ...
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	var err error
	c.config.dsnCfg, err = ParseDSN(c.config.DSN)
	if err == nil {
		c.logger.Info(ecode.MsgClientMysqlOpenStart, elog.FieldMod("gorm"), elog.FieldAddr(c.config.dsnCfg.Addr), elog.FieldName(c.config.dsnCfg.DBName))
	} else {
		c.logger.Panic(ecode.MsgClientMysqlOpenStart, elog.FieldMod("gorm"), elog.FieldErr(err))
	}

	if c.config.Debug {
		c.WithInterceptor(debugInterceptor)
	}
	if !c.config.DisableTrace {
		c.WithInterceptor(traceInterceptor)
	}

	if !c.config.DisableMetric {
		c.WithInterceptor(metricInterceptor)
	}

	component, err := newComponent(c.name, c.config, c.logger)
	if err != nil {
		if c.config.OnDialError == "panic" {
			c.logger.Panic("open mysql", elog.FieldMod("gorm"), elog.FieldErrKind(ecode.ErrKindRequestErr), elog.FieldErr(err), elog.FieldAddr(c.config.dsnCfg.Addr), elog.FieldValueAny(c.config))
		} else {
			metric.LibHandleCounter.Inc(metric.TypeGorm, c.name+".ping", c.config.dsnCfg.Addr, "open err")
			c.logger.Error("open mysql", elog.FieldMod("gorm"), elog.FieldErrKind(ecode.ErrKindRequestErr), elog.FieldErr(err), elog.FieldAddr(c.config.dsnCfg.Addr), elog.FieldValueAny(c.config))
			return component
		}
	}

	if err := component.DB.DB().Ping(); err != nil {
		c.logger.Panic("ping mysql", elog.FieldMod("gorm"), elog.FieldErrKind(ecode.ErrKindRequestErr), elog.FieldErr(err), elog.FieldValueAny(c.config))
	}

	// store db
	instances.Store(c.name, component)
	return component
}

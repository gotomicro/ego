package ejob

import (
	"context"
	"time"

	"github.com/gotomicro/ego/core/eflag"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/etrace"
	"github.com/gotomicro/ego/core/standard"
)

func init() {
	eflag.Register(
		&eflag.StringFlag{
			Name:    "job",
			Usage:   "--job",
			Default: "",
		},
	)
}

const PackageName = "task.ejob"

// Component ...
type Component struct {
	name   string
	config *Config
	logger *elog.Component
}

func newComponent(name string, config *Config, logger *elog.Component) *Component {
	return &Component{
		name:   name,
		config: config,
		logger: logger,
	}
}

func (c *Component) Name() string {
	return c.config.Name
}

func (c *Component) PackageName() string {
	return PackageName
}

func (c *Component) Init() error {
	return nil
}

func (c *Component) Start() error {
	span, ctx := etrace.StartSpanFromContext(
		context.Background(),
		"ego-job",
	)
	defer span.Finish()
	traceId := etrace.ExtractTraceID(ctx)
	beg := time.Now()
	c.logger.Info("start ejob", elog.FieldName(c.name), elog.FieldTid(traceId))
	err := c.config.startFunc(ctx)
	if err != nil {
		c.logger.Error("stop ejob", elog.FieldName(c.name), elog.FieldErr(err), elog.FieldCost(time.Since(beg)), elog.FieldTid(traceId))
	} else {
		c.logger.Info("stop ejob", elog.FieldName(c.name), elog.FieldCost(time.Since(beg)), elog.FieldTid(traceId))
	}
	return err
}

func (c *Component) Stop() error {
	return nil
}

// Ejob ...
type Ejob interface {
	standard.Component
}

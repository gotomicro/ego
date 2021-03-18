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

// PackageName 包名
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

// Name 配置名称
func (c *Component) Name() string {
	return c.config.Name
}

// PackageName 包名
func (c *Component) PackageName() string {
	return PackageName
}

// Init 初始化
func (c *Component) Init() error {
	return nil
}

// Start 启动
func (c *Component) Start() error {
	span, ctx := etrace.StartSpanFromContext(
		context.Background(),
		"ego-job",
	)
	defer span.Finish()
	traceID := etrace.ExtractTraceID(ctx)
	beg := time.Now()
	c.logger.Info("start ejob", elog.FieldName(c.name), elog.FieldTid(traceID))
	err := c.config.startFunc(ctx)
	if err != nil {
		c.logger.Error("stop ejob", elog.FieldName(c.name), elog.FieldErr(err), elog.FieldCost(time.Since(beg)), elog.FieldTid(traceID))
	} else {
		c.logger.Info("stop ejob", elog.FieldName(c.name), elog.FieldCost(time.Since(beg)), elog.FieldTid(traceID))
	}
	return err
}

// Stop ...
func (c *Component) Stop() error {
	return nil
}

// Ejob ...
type Ejob interface {
	standard.Component
}

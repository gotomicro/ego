package ejob

import (
	"github.com/gotomicro/ego/core/eflag"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/standard"
	"time"
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
	beg := time.Now()
	c.logger.Info("start ejob", elog.FieldName(c.name))
	err := c.config.startFunc()
	if err != nil {
		c.logger.Error("stop ejob", elog.FieldName(c.name), elog.FieldErr(err), elog.FieldCost(time.Since(beg)))
	} else {
		c.logger.Info("stop ejob", elog.FieldName(c.name), elog.FieldCost(time.Since(beg)))
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

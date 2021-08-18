package ejob

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

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

// Context Job Context
type Context struct {
	Ctx     context.Context
	Writer  http.ResponseWriter
	Request *http.Request
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

func (c *Component) trace(ctx context.Context) {
	var (
		traceID = etrace.ExtractTraceID(ctx)
		fields  = []elog.Field{elog.FieldName(c.name)}
		err     error
	)

	// 如果设置了链路，增加链路信息
	if opentracing.IsGlobalTracerRegistered() {
		fields = append(fields, elog.FieldTid(traceID))
	}
	beg := time.Now()
	c.logger.Info("start ejob", fields...)
	defer func() {
		if rec := recover(); rec != nil {
			switch rec := rec.(type) {
			case error:
				err = rec
			default:
				err = fmt.Errorf("%v", rec)
			}

			stack := make([]byte, 4096)
			length := runtime.Stack(stack, true)
			fields = append(fields, zap.ByteString("stack", stack[:length]))
		}
		if err != nil {
			fields = append(fields, elog.FieldErr(err), elog.FieldCost(time.Since(beg)))
			c.logger.Error("start  ejob", fields...)
		} else {
			fields = append(fields, elog.FieldCost(time.Since(beg)))
			c.logger.Info("start  ejob", fields...)
		}
	}()
}

// StartHTTP ...
func (c *Component) StartHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	span, ctx := etrace.StartSpanFromContext(
		context.Background(),
		"ego-job",
	)
	defer span.Finish()
	c.trace(ctx)
	return c.config.startFunc(Context{
		Ctx:     ctx,
		Writer:  w,
		Request: r,
	})
}

// Start 启动
func (c *Component) Start() (err error) {
	span, ctx := etrace.StartSpanFromContext(
		context.Background(),
		"ego-job",
	)
	defer span.Finish()
	c.trace(ctx)
	return c.config.startFunc(Context{
		Ctx: ctx,
	})
}

// Stop ...
func (c *Component) Stop() error {
	return nil
}

// Ejob ...
type Ejob interface {
	standard.Component
}

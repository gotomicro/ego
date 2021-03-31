package ecron

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/emetric"
	"github.com/gotomicro/ego/core/etrace"
)

type wrappedJob struct {
	NamedJob
	logger *elog.Component
}

// Run ...
func (wj wrappedJob) Run() {
	wj.run()
}

func (wj wrappedJob) run() {
	span, ctx := etrace.StartSpanFromContext(
		context.Background(),
		"ego-cron",
	)
	defer span.Finish()
	traceId := etrace.ExtractTraceID(ctx)
	emetric.JobHandleCounter.Inc("cron", wj.Name(), "begin")
	var fields = []elog.Field{zap.String("name", wj.Name())}
	// 如果设置了链路，增加链路信息
	if opentracing.IsGlobalTracerRegistered() {
		fields = append(fields, elog.FieldTid(traceId))
	}

	wj.logger.Info("cron start", fields...)
	var beg = time.Now()
	defer func() {
		var err error
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
			fields = append(fields, elog.FieldErr(err), elog.Duration("cost", time.Since(beg)))
			wj.logger.Error("cron end", fields...)
		} else {
			wj.logger.Info("cron end", fields...)
		}
		emetric.JobHandleHistogram.Observe(time.Since(beg).Seconds(), "cron", wj.Name())
	}()

	err := wj.NamedJob.Run(ctx)
	if err != nil {
		fields = append(fields, elog.FieldErr(err))
		wj.logger.Error("cron run failed", fields...)
	}
}

package main

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/etrace"
	"github.com/gotomicro/ego/task/ejob"
)

// export EGO_DEBUG=true && go run main.go --job=jobrunner  --config=config.toml
func main() {
	if err := ego.New().Job(
		ejob.Job("jobrunner", runner),
		ejob.Job("job1", job1),
	).Run(); err != nil {
		elog.Error("start up", zap.Error(err))
	}
}

func runner(ctx ejob.Context) error {
	fmt.Println("i am job runner, traceId: ", etrace.ExtractTraceID(ctx.Ctx))
	return errors.New("i am error")
}

func job1(ctx ejob.Context) error {
	fmt.Println("i am job runner, traceId: ", etrace.ExtractTraceID(ctx.Ctx))
	return errors.New("i am error")
}

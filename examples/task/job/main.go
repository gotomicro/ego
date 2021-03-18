package main

import (
	"context"
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
	if err := ego.New().Job(NewJobRunner()).Run(); err != nil {
		elog.Error("start up", zap.Error(err))
	}
}

// NewJobRunner 创建新的job
func NewJobRunner() *ejob.Component {
	return ejob.DefaultContainer().Build(
		ejob.WithName("jobrunner"),
		ejob.WithStartFunc(runner),
	)
}

func runner(ctx context.Context) error {
	fmt.Println("i am job runner, traceId: ", etrace.ExtractTraceID(ctx))
	return errors.New("i am error")
}

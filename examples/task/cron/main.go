package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/etrace"
	"github.com/gotomicro/ego/task/ecron"
)

//  export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
	err := ego.New().Cron(cronJob1(), cronJob2()).Run()
	if err != nil {
		elog.Panic("startup engine", elog.FieldErr(err))
	}
}

// 异常任务
func cronJob1() ecron.Ecron {
	job := func(ctx context.Context) error {
		elog.Info("info job1", elog.FieldTid(etrace.ExtractTraceID(ctx)))
		elog.Warn("warn job1", elog.FieldTid(etrace.ExtractTraceID(ctx)))
		fmt.Println("run job1", elog.FieldTid(etrace.ExtractTraceID(ctx)))
		return errors.New("exec job1 error")
	}

	cron := ecron.Load("cron.test").Build(ecron.WithJob(job))
	return cron
}

// 正常任务
func cronJob2() ecron.Ecron {
	job := func(ctx context.Context) error {
		elog.Info("info job2", elog.FieldTid(etrace.ExtractTraceID(ctx)))
		elog.Warn("warn job2", elog.FieldTid(etrace.ExtractTraceID(ctx)))
		fmt.Println("run job2", elog.FieldTid(etrace.ExtractTraceID(ctx)))
		return nil
	}

	cron := ecron.Load("cron.test").Build(ecron.WithJob(job))
	return cron
}

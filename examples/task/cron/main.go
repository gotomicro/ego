package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/etrace"
	"github.com/gotomicro/ego/task/ecron"
)

//  export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
	err := ego.New().Cron(cron1()).Run()
	if err != nil {
		elog.Panic("startup engine", elog.FieldErr(err))
	}
}

func cron1() ecron.Ecron {
	cron := ecron.Load("cron.test").Build()
	cron.Schedule(ecron.Every(time.Second*10), ecron.FuncJob(execJob))
	cron.Schedule(ecron.Every(time.Second*10), ecron.FuncJob(execJob2))
	return cron
}

// 异常任务
func execJob(ctx context.Context) error {
	elog.Info("info job", elog.FieldTid(etrace.ExtractTraceID(ctx)))
	elog.Warn("warn job", elog.FieldTid(etrace.ExtractTraceID(ctx)))
	fmt.Println("run job", elog.FieldTid(etrace.ExtractTraceID(ctx)))
	return errors.New("exec job1 error")
}

// 正常任务
func execJob2(ctx context.Context) error {
	elog.Info("info job2", elog.FieldTid(etrace.ExtractTraceID(ctx)))
	elog.Warn("warn job2", elog.FieldTid(etrace.ExtractTraceID(ctx)))
	fmt.Println("run job2", elog.FieldTid(etrace.ExtractTraceID(ctx)))
	return nil
}

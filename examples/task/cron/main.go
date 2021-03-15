package main

import (
	"errors"
	"fmt"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/task/ecron"
	"time"
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
func execJob() error {
	elog.Info("info job")
	elog.Warn("warn job")
	fmt.Println("run job")
	return errors.New("exec job1 error")
}

// 正常任务
func execJob2() error {
	elog.Info("info job2")
	elog.Warn("warn job2")
	fmt.Println("run job2")
	return nil
}

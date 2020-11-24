package main

import (
	"errors"
	"fmt"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"go.uber.org/zap"
)

// export EGO_DEBUG=true && go run main.go --job=jobrunner
func main() {
	if err := ego.New().Job(NewJobRunner()).Run(); err != nil {
		elog.Error("start up", zap.Error(err))
	}
}

type JobRunner struct {
	JobName string
}

func NewJobRunner() *JobRunner {
	return &JobRunner{
		JobName: "jobrunner",
	}
}

func (j *JobRunner) Run() error {
	fmt.Println("i am job runner")
	return errors.New("i am error")
}

func (j *JobRunner) GetJobName() string {
	return j.JobName
}

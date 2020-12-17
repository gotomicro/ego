package main

import (
	"errors"
	"fmt"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/task/ejob"
	"go.uber.org/zap"
)

// export EGO_DEBUG=true && go run main.go --job=jobrunner
func main() {
	if err := ego.New().Job(NewJobRunner()).Run(); err != nil {
		elog.Error("start up", zap.Error(err))
	}
}

func NewJobRunner() *ejob.Component {
	return ejob.DefaultContainer().Build(
		ejob.WithName("jobrunner"),
		ejob.WithStartFunc(runner),
	)
}

func runner() error {
	fmt.Println("i am job runner")
	return errors.New("i am error")
}

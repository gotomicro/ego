package main

import (
	"errors"

	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
)

func main() {
	err := ego.New().Invoker(func() error {
		elog.Info("logger info", elog.String("gopher", "ego"), elog.String("type", "command"))
		return errors.New("i am panic")
	}).Run()
	if err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

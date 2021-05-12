package main

import (
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
)

//  export EGO_DEBUG=true && go run main.go
func main() {
	err := ego.New().Invoker(func() error {
		elog.Info("logger info", elog.String("gopher", "ego"), elog.String("type", "command"), elog.Any("aaa", map[string]interface{}{"aa": "bb"}))
		return nil
	}).Run()
	if err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

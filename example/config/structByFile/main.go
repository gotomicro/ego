package main

import (
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/conf"
	"github.com/gotomicro/ego/core/elog"
)

//  export EGO_DEBUG=true && go run main.go --config=config.toml --watch=false
func main() {
	err := ego.New(func() error {
		p := People{}
		err := conf.UnmarshalKey("people", &p)
		if err != nil {
			panic(err.Error())
		}
		elog.Info("people info", elog.String("name", p.Name), elog.String("type", "structByFile"))
		return nil
	}).Run()
	if err != nil {
		elog.Panic("startup", elog.Any("err", err))
	}
}

type People struct {
	Name string
}

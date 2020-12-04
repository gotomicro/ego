package main

import (
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
)

// export EGO_DEBUG=true && go run main.go  --config=config.toml --watch=false
func main() {
	if err := ego.New().Invoker(func() error {
		peopleName := econf.GetString("people.name")
		elog.Info("people info", elog.String("name", peopleName), elog.String("type", "onelineByFile"))
		return nil
	}).Run(); err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

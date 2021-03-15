package main

import (
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"time"
)

//  export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
	err := ego.New(ego.WithHang(true)).Invoker(func() error {
		go func() {
			for {
				elog.Info("logger info", elog.String("gopher", "ego1"), elog.String("type", "file"))
				elog.Debug("logger debug", elog.String("gopher", "ego2"), elog.String("type", "file"))
				time.Sleep(1 * time.Second)
			}
		}()
		return nil
	}).Run()
	if err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

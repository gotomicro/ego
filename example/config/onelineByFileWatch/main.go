package main

import (
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/conf"
	"github.com/gotomicro/ego/core/elog"
	"time"
)

//  export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
	if err := ego.New(func() error {
		go func() {
			// 循环打印配置
			for {
				time.Sleep(3 * time.Second)
				peopleName := conf.GetString("people.name")
				elog.Info("people info", elog.String("name", peopleName), elog.String("type", "onelineByFileWatch"))
			}
		}()
		return nil
	}).Hang(true).Run(); err != nil {
		elog.Panic("startup", elog.Any("err", err))
	}
}

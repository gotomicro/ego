package main

import (
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/conf"
	"github.com/gotomicro/ego/core/elog"
	"time"
)

//  export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
	if err := ego.New(ego.WithHang(true)).Invoker(func() error {
		p := People{}
		// 初始化
		err := conf.UnmarshalKey("people", &p)
		if err != nil {
			panic(err.Error())
		}
		// 监听
		conf.OnChange(func(config *conf.Configuration) {
			err := config.UnmarshalKey("people", &p)
			if err != nil {
				panic(err.Error())
			}
		})

		go func() {
			// 循环打印配置
			for {
				time.Sleep(1 * time.Second)
				elog.Info("people info", elog.String("name", p.Name), elog.String("type", "structByFileWatch"))
			}
		}()
		return nil
	}).Run(); err != nil {
		elog.Panic("startup", elog.Any("err", err))
	}
}

type People struct {
	Name string
}

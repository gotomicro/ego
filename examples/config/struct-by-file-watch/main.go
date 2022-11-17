package main

import (
	"time"

	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
)

//  export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
	if err := ego.New(ego.WithHang(true)).Invoker(func() error {
		p := People{}
		// 初始化
		err := econf.UnmarshalKey("people", &p)
		if err != nil {
			panic(err.Error())
		}
		// 监听
		econf.OnChange(func(config *econf.Configuration) {
			err := config.UnmarshalKey("people", &p)
			if err != nil {
				elog.Panic("unmarshal", elog.FieldErr(err))
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
		elog.Panic("startup", elog.FieldErr(err))
	}
}

// People ...
type People struct {
	Name string
}

package xdebug

import (
	"fmt"
	"github.com/gotomicro/ego/core/util/xcolor"
	"time"
)

// 配置名、目标地址、耗时、请求参数、响应数据
func Info(compName string, addr string, cost time.Duration, req interface{}, reply interface{}) {
	fmt.Printf("%s %s %s %s %s => %s\n", time.Now().Format("2006-01-02 15:04:05.000"), xcolor.Green(compName), xcolor.Green(addr), xcolor.Yellow(fmt.Sprintf("[%vms]", float64(cost.Microseconds())/1000)), xcolor.Blue(fmt.Sprintf("%v", req)), xcolor.Blue(fmt.Sprintf("%v", reply)))
}

func Error(compName string, addr string, cost time.Duration, req string, err string) {
	fmt.Printf("%s %s %s %s %s => %s\n", time.Now().Format("2006-01-02 15:04:05.000"), xcolor.Red(compName), xcolor.Red(addr), xcolor.Yellow(fmt.Sprintf("[%vms]", float64(cost.Microseconds())/1000)), xcolor.Blue(fmt.Sprintf("%v", req)), xcolor.Red(err))
}

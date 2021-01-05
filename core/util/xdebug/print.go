package xdebug

import (
	"fmt"
	"time"

	"github.com/gotomicro/ego/core/util/xcolor"
)

// MakeReqResInfo 以info级别打印配置名、目标地址、耗时、请求数据、响应数据
func MakeReqResInfo(compName string, addr string, cost time.Duration, req interface{}, reply interface{}) string {
	return fmt.Sprintf("%s %s %s %s => %s\n", xcolor.Green(compName), xcolor.Green(addr), xcolor.Yellow(fmt.Sprintf("[%vms]", float64(cost.Microseconds())/1000)), xcolor.Blue(fmt.Sprintf("%v", req)), xcolor.Blue(fmt.Sprintf("%v", reply)))
}

// MakeReqResInfo 以error级别打印配置名、目标地址、耗时、请求数据、响应数据
func MakeReqResError(compName string, addr string, cost time.Duration, req string, err string) string {
	return fmt.Sprintf("%s %s %s %s => %s\n", xcolor.Red(compName), xcolor.Red(addr), xcolor.Yellow(fmt.Sprintf("[%vms]", float64(cost.Microseconds())/1000)), xcolor.Blue(fmt.Sprintf("%v", req)), xcolor.Red(err))
}

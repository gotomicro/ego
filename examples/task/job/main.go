package main

import (
	"errors"
	"fmt"
	"io/ioutil"

	"go.uber.org/zap"

	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/etrace"
	"github.com/gotomicro/ego/task/ejob"
)

// export EGO_DEBUG=true && go run main.go --job=jobrunner  --config=config.toml
func main() {
	if err := ego.New().Job(
		ejob.Job("jobrunner", runner),
		ejob.Job("job1", job1),
	).Run(); err != nil {
		elog.Error("start up", zap.Error(err))
	}
}

type data struct {
	Test int
}

// 测试job链接  export EGO_DEBUG=true && go run main.go --config=config.toml --job=jobrunner --job-data='{"test":1}' --job-header='test=3;asdf=4'
// 支持job里存入data数据，和http请求保持统一
func runner(ctx ejob.Context) error {
	str, _ := ioutil.ReadAll(ctx.Request.Body)
	fmt.Printf("str--------------->"+"%+v\n", string(str))
	fmt.Printf("str--------------->"+"%+v\n", ctx.Request.Header.Get("test"))
	fmt.Printf("ctx.Request.URL--------------->%s", ctx.Request)
	fmt.Println("i am job runner, traceId: ", etrace.ExtractTraceID(ctx.Ctx))
	return errors.New("i am error")
}

func job1(ctx ejob.Context) error {
	fmt.Println("i am job runner, traceId: ", etrace.ExtractTraceID(ctx.Ctx))
	return errors.New("i am error")
}

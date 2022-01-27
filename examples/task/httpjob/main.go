package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egovernor"
	"github.com/gotomicro/ego/task/ejob"
	"go.uber.org/zap"
)

// 如果是Job 命令行执行  export EGO_DEBUG=true && go run main.go --config=config.toml --job=job --job-data='{"username":"ego"}' --job-header='test=1'
// 如果是Job HTTP执行  1 export EGO_DEBUG=true && go run main.go --config=config.toml
// 如果是Job HTTP执行  2 curl -v -XPOST -d '{"username":"ego"}' -H 'X-Ego-Job-Name:job' -H 'X-Ego-Job-RunID:xxxx' -H 'test=1' http://127.0.0.1:9003/jobs
func main() {
	if err := ego.New().Job(
		ejob.Job("job", job),
	).Serve(
		egovernor.Load("server.governor").Build(),
	).Run(); err != nil {
		elog.Error("start up", zap.Error(err))
	}
}

type data struct {
	Username string
}

func job(ctx ejob.Context) error {
	bytes, _ := ioutil.ReadAll(ctx.Request.Body)
	d := data{}
	_ = json.Unmarshal(bytes, &d)
	fmt.Println(d.Username)
	fmt.Println(ctx.Request.Header.Get("test"))
	ctx.Writer.Write([]byte("i am ok"))
	return nil
}

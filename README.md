## EGO
[![Go](https://github.com/gotomicro/ego/workflows/Go/badge.svg?branch=master)](https://github.com/gotomicro/ego/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/gotomicro/ego)](https://goreportcard.com/report/github.com/gotomicro/ego)
[![codecov](https://codecov.io/gh/gotomicro/ego/branch/master/graph/badge.svg)](https://codecov.io/gh/gotomicro/ego)
[![goproxy.cn](https://goproxy.cn/stats/github.com/gotomicro/ego/badges/download-count.svg)](https://goproxy.cn/stats/github.com/gotomicro/ego)
[![Release](https://img.shields.io/github/v/release/gotomicro/ego.svg?style=flat-square)](https://github.com/gotomicro/ego)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 帮助文档
[https://ego.gocn.vip](https://ego.gocn.vip)

## Quick Start

### HelloWorld
```package main
import (
   "github.com/gin-gonic/gin"
   "github.com/gotomicro/ego"
   "github.com/gotomicro/ego/core/elog"
   "github.com/gotomicro/ego/server"
   "github.com/gotomicro/ego/server/egin"
)
//  export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
   if err := ego.New().Serve(func() *egin.Component {
      server := egin.Load("server.http").Build()
      server.GET("/hello", func(ctx *gin.Context) {
         ctx.JSON(200, "Hello EGO")
         return
      })
      return server
   }()).Run(); err != nil {
      elog.Panic("startup", elog.FieldErr(err))
   }
}
```

### 使用命令行运行
```
export EGO_DEBUG=true # 默认日志输出到logs目录，开启dev后日志输出到终端
go run main.go --config=config.toml
```

### 如下所示
![图片](./docs/images/startup.png)


这个时候我们可以发送一个指令，得到如下结果
```
➜  helloworld git:(master) ✗ curl http://127.0.0.1:9001/hello
"Hello Ego"%  
```

### 更加友好的包编译

使用scripts文件夹里的[包编译](./example/build)，可以看到优雅的version提示

![图片](./docs/images/version.png)
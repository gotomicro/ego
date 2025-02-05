package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egin"
	"github.com/gotomicro/ego/server/egovernor"
)

var (
	// server   *http.Server
	// listener net.Listener = nil

	reload  = flag.Bool("reload", false, "listen on fd open 3 (internal use only)")
	message = flag.String("message", "Hello World", "message to send")
)

// export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
	if err := ego.New().ReloadServe(func() *egin.Component {
		var server *egin.Component
		elog.Info("reload", elog.Any("reload", *reload))
		if *reload {
			f := os.NewFile(3, "")
			listener, err := net.FileListener(f)
			if err != nil {
				elog.Panic("listener error", elog.FieldErr(err))
			}
			server = egin.Load("server.http").Build(egin.WithListener(listener))
		} else {
			server = egin.Load("server.http").Build()
		}
		server.GET("/hello", func(ctx *gin.Context) {
			ctx.JSON(200, fmt.Sprintf("Hello Check %s, PID: %d", *message, os.Getpid()))
		})
		server.GET("/kill", func(ctx *gin.Context) {
			// 热重启
			pid := os.Getpid()
			elog.Info("kill", elog.Int("pid", pid))
			err := syscall.Kill(pid, syscall.SIGQUIT)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			}
			ctx.JSON(200, "success")
		})
		server.GET("/reload", func(ctx *gin.Context) {
			// 热重启
			pid := os.Getpid()
			elog.Info("reload", elog.Int("pid", pid))
			err := syscall.Kill(pid, syscall.SIGUSR1)
			if err != nil {
				ctx.JSON(500, err.Error())
				return
			}
			ctx.JSON(200, "success")
		})
		return server
	}()).Serve(egovernor.Load("server.governor").Build()).Run(); err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

//
// func main2() {
// 	var err error
//
// 	// 解析参数
// 	flag.Parse()
//
// 	http.HandleFunc("/test", handler)
// 	server = &http.Server{Addr: ":3000"}
//
// 	// 设置监听器的监听对象（新建的或已存在的 socket 描述符）
// 	if *graceful {
// 		// 子进程监听父进程传递的 socket 描述符
// 		log.Println("listening on the existing file descriptor 3")
// 		// 子进程的 0, 1, 2 是预留给标准输入、标准输出、错误输出，故传递的 socket 描述符
// 		// 应放在子进程的 3
// 		f := os.NewFile(3, "")
// 		listener, err = net.FileListener(f)
// 	} else {
// 		// 父进程监听新建的 socket 描述符
// 		log.Println("listening on a new file descriptor")
// 		listener, err = net.Listen("tcp", server.Addr)
// 	}
// 	if err != nil {
// 		log.Fatalf("listener error: %v", err)
// 	}
//
// 	go func() {
// 		err = server.Serve(listener)
// 		log.Printf("server.Serve err: %v\n", err)
// 	}()
// 	// 监听信号
// 	handleSignal()
// 	log.Println("signal end")
// }
//
// func handleSignal() {
// 	ch := make(chan os.Signal, 1)
// 	// 监听信号
// 	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
// 	for {
// 		sig := <-ch
// 		log.Printf("signal receive: %v\n", sig)
// 		ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
// 		switch sig {
// 		case syscall.SIGINT, syscall.SIGTERM: // 终止进程执行
// 			log.Println("shutdown")
// 			signal.Stop(ch)
// 			server.Shutdown(ctx)
// 			log.Println("graceful shutdown")
// 			return
// 		case syscall.SIGUSR2: // 进程热重启
// 			log.Println("reload")
// 			err := reload() // 执行热重启函数
// 			if err != nil {
// 				log.Fatalf("graceful reload error: %v", err)
// 			}
// 			server.Shutdown(ctx)
// 			log.Println("graceful reload")
// 			return
// 		}
// 	}
// }
//
// func reload() error {
// 	tl, ok := listener.(*net.TCPListener)
// 	if !ok {
// 		return errors.New("listener is not tcp listener")
// 	}
// 	// 获取 socket 描述符
// 	f, err := tl.File()
// 	if err != nil {
// 		return err
// 	}
// 	// 设置传递给子进程的参数（包含 socket 描述符）
// 	args := []string{"-graceful"}
// 	cmd := exec.Command(os.Args[0], args...)
// 	cmd.Stdout = os.Stdout         // 标准输出
// 	cmd.Stderr = os.Stderr         // 错误输出
// 	cmd.ExtraFiles = []*os.File{f} // 文件描述符
// 	// 新建并执行子进程
// 	return cmd.Start()
// }

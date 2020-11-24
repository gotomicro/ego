package signals

import (
	"os"
	"syscall"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func kill(sig os.Signal) {
	pro, _ := os.FindProcess(os.Getpid())
	pro.Signal(sig)
}
func TestShutdownSIGQUIT(t *testing.T) {
	quit := make(chan struct{})
	Convey("test shutdown signal by SIGQUIT", t, func(c C) {
		fn := func(grace bool) {
			c.So(grace, ShouldEqual, false)
			close(quit)
		}
		Shutdown(fn)
		kill(syscall.SIGQUIT)
		<-quit
	})
}

// func TestShutdownSIGINT(t *testing.T) {
// 	quit := make(chan struct{})
// 	Convey("test shutdown signal by SIGINT", t, func(c C) {
// 		fn := func(grace bool) {
// 			c.So(grace, ShouldEqual, true)
// 			close(quit)
// 		}
// 		Shutdown(fn)
// 		kill(syscall.SIGINT)
// 		<-quit
// 	})
// }

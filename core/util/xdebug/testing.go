package xdebug

import (
	"flag"
	"fmt"
	"github.com/gotomicro/ego/core/app"
	"sync"

	"github.com/gotomicro/ego/core/util/xcolor"
	"github.com/gotomicro/ego/core/util/xstring"
	"github.com/tidwall/pretty"
)

var (
	isTestingMode bool
	//isDevelopmentMode = os.Getenv("EGO_MODE") == "dev"
)

//func init() {
//	if isDevelopmentMode {
//		xlog.DefaultLogger.SetLevel(xlog.DebugLevel)
//		xlog.EgoLogger.SetLevel(xlog.DebugLevel)
//	}
//}

// IsTestingMode 判断是否在测试模式下
var onceTest = sync.Once{}

// IsTestingMode ...
func IsTestingMode() bool {
	onceTest.Do(func() {
		isTestingMode = flag.Lookup("test.v") != nil
	})

	return isTestingMode
}

// IsDevelopmentMode 判断是否是生产模式
func IsDevelopmentMode() bool {
	//return isDevelopmentMode || isTestingMode
	return app.IsDevelopmentMode()
}

// IfPanic ...
func IfPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// PrettyJsonPrint ...
func PrettyJsonPrint(message string, obj interface{}) {
	if !IsDevelopmentMode() {
		return
	}
	fmt.Printf("%s => %s\n",
		xcolor.Red(message),
		pretty.Color(
			pretty.Pretty([]byte(xstring.PrettyJson(obj))),
			pretty.TerminalStyle,
		),
	)
}

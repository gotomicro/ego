package xdebug

import (
	"fmt"

	"github.com/gotomicro/ego/core/util/xcolor"
	"github.com/gotomicro/ego/core/util/xstring"
	"github.com/tidwall/pretty"
)

// DebugObject ...
func PrintObject(message string, obj interface{}) {
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

// DebugBytes ...
func DebugBytes(obj interface{}) string {
	return string(pretty.Color(pretty.Pretty([]byte(xstring.Json(obj))), pretty.TerminalStyle))
}

// PrintKV ...
func PrintKV(key string, val string) {
	if !IsDevelopmentMode() {
		return
	}
	fmt.Printf("%-50s => %s\n", xcolor.Red(key), xcolor.Green(val))
}

// PrettyKVWithPrefix ...
func PrintKVWithPrefix(prefix string, key string, val string) {
	if !IsDevelopmentMode() {
		return
	}
	fmt.Printf("%-8s]> %-30s => %s\n", prefix, xcolor.Red(key), xcolor.Blue(val))
}

// PrintMap ...
func PrintMap(data map[string]interface{}) {
	if !IsDevelopmentMode() {
		return
	}
	for key, val := range data {
		fmt.Printf("%-20s : %s\n", xcolor.Red(key), fmt.Sprintf("%+v", val))
	}
}

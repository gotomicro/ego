package xstring

import (
	"reflect"
	"runtime"
)

// FunctionName ...
func FunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// ObjectName ...
func ObjectName(i interface{}) string {
	typ := reflect.TypeOf(i)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return typ.PkgPath() + "." + typ.Name()
}

// CallerName ...
func CallerName(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	return runtime.FuncForPC(pc).Name()
}

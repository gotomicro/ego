package xstring

import (
	"reflect"
	"runtime"
)

// FunctionName returns the Function's name of given pointer.
func FunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// ObjectName returns the object's path and name of given pointer.
// Deprecated: this function will be moved to internal package, user should not use it any more.
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

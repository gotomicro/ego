package tools

import (
	"context"
	"fmt"
	"go/format"
	"log"
	"reflect"

	"github.com/spf13/cast"
	"google.golang.org/grpc/metadata"
)

// GrpcHeaderValue 获取context value
func GrpcHeaderValue(ctx context.Context, key string) string {
	if key == "" {
		return ""
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	if len(md.Get(key)) > 0 {
		return md.Get(key)[0]
	}
	return ""
}

// ContextValue gRPC日志获取context value
func ContextValue(ctx context.Context, key string) string {
	if key == "" {
		return ""
	}
	return cast.ToString(ctx.Value(key))
}

// ToSliceStringMap casts an empty interface to []map[string]interface{} ignoring error
func ToSliceStringMap(i interface{}) []map[string]interface{} {
	var s = make([]map[string]interface{}, 0)
	switch v := i.(type) {
	case []interface{}:
		for _, u := range v {
			s = append(s, cast.ToStringMap(u))
		}
		return s
	case []map[string]interface{}:
		s = append(s, v...)
		return s
	default:
		log.Printf("unable to Cast %#v of type %v to []map[string]interface{}", i, reflect.TypeOf(i))
		return s
	}
}

// GoFmt 格式化代码
func GoFmt(buf []byte) []byte {
	formatted, err := format.Source(buf)
	if err != nil {
		panic(fmt.Errorf("%s\nOriginal code:\n%s", err.Error(), buf))
	}
	return formatted
}

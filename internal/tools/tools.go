package tools

import (
	"context"
	"fmt"
	"go/format"
	"log"
	"reflect"
	"strings"

	"github.com/spf13/cast"
	"google.golang.org/grpc/metadata"

	"github.com/gotomicro/ego/core/transport"
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
	// 小写
	return strings.Join(md.Get(key), ";")
}

// ContextValue gRPC日志获取context value
func ContextValue(ctx context.Context, key string) string {
	if key == "" {
		return ""
	}
	return cast.ToString(transport.Value(ctx, key))
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

// GoFmt 格式化Go
func GoFmt(buf []byte) []byte {
	formatted, err := format.Source(buf)
	if err != nil {
		panic(fmt.Errorf("%s\nOriginal code:\n%s", err.Error(), buf))
	}
	return formatted
}

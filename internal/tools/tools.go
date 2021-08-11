package tools

import (
	"context"
	"fmt"
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
	v, _ := toSliceStringMapE(i)
	return v
}

func toSliceStringMapE(i interface{}) ([]map[string]interface{}, error) {
	var s = make([]map[string]interface{}, 0)

	switch v := i.(type) {
	case []interface{}:
		for _, u := range v {
			s = append(s, cast.ToStringMap(u))
		}
		return s, nil
	case []map[string]interface{}:
		s = append(s, v...)
		return s, nil
	default:
		return s, fmt.Errorf("unable to Cast %#v of type %v to []map[string]interface{}", i, reflect.TypeOf(i))
	}
}

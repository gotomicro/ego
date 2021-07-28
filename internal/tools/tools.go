package tools

import (
	"context"
	"strings"

	"github.com/gotomicro/ego/core/transport"
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

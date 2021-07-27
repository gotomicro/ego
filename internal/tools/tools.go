package tools

import (
	"context"
	"strings"

	"github.com/spf13/cast"
	"google.golang.org/grpc/metadata"
)

// GetContextValue 获取context value
func GetContextValue(ctx context.Context, key string) string {
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

// LoggerGrpcContextValue gRPC日志获取context value
func LoggerGrpcContextValue(ctx context.Context, key string) string {
	value := GetContextValue(ctx, key)
	if value != "" {
		return value
	}
	return cast.ToString(ctx.Value(key))
}

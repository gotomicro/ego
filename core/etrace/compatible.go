package etrace

import (
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
)

// CompatibleExtractHttpTraceId ...
// Deprecated 该方法会在v0.9.0移除
func CompatibleExtractHttpTraceId(header http.Header) {
	xTraceId := header.Get("X-Trace-Id")
	if xTraceId != "" {
		header.Set("Traceparent", CompatibleParse(xTraceId))
	}
}

// CompatibleExtractGrpcTraceId ...
// Deprecated 该方法会在v0.9.0移除
func CompatibleExtractGrpcTraceId(header metadata.MD) {
	xTraceId := header.Get("x-trace-id")
	fmt.Println(len(xTraceId))
	if len(xTraceId) > 0 {
		header.Set("Traceparent", CompatibleParse(xTraceId[0]))
	}
}

// opentrace: 18af9db18a77f4b7:18af9db18a77f4b7:0000000000000000:0
// opentelemetry: 00-18af9db18a77f4b718af9db18a77f4b7-18af9db18a77f4b7-00
// https://www.w3.org/TR/trace-context/
func CompatibleParse(traceId string) string {
	traceArr := strings.Split(traceId, ":")
	if len(traceArr) == 4 {
		return "00-" + traceArr[0] + traceArr[1] + "-" + traceArr[1] + "-0" + traceArr[3]
	}
	return ""
}

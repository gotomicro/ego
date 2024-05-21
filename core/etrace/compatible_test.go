package etrace

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestCompatibleExtractHTTPTraceID(t *testing.T) {
	header := make(http.Header)
	header.Set("X-Trace-Id", "111222")
	CompatibleExtractHTTPTraceID(header)
	var tp = header.Get("X-Trace-Id")
	assert.Equal(t, "111222", tp)
}

func TestCompatibleExtractGrpcTraceID(t *testing.T) {
	header := make(metadata.MD)
	CompatibleExtractGrpcTraceID(header)
	assert.NoError(t, nil)
}

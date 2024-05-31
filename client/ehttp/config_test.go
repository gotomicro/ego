package ehttp

import (
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/util/xtime"
)

func Test_DefaultConfig(t *testing.T) {
	assert.True(t, reflect.DeepEqual(&Config{
		Addr:                       "",
		Debug:                      false,
		RawDebug:                   false,
		ReadTimeout:                xtime.Duration("2s"),
		SlowLogThreshold:           xtime.Duration("500ms"),
		IdleConnTimeout:            90 * time.Second,
		MaxIdleConns:               100,
		MaxIdleConnsPerHost:        runtime.GOMAXPROCS(0) + 1,
		EnableTraceInterceptor:     true,
		EnableKeepAlives:           true,
		EnableAccessInterceptor:    false,
		EnableAccessInterceptorRes: false,
		EnableMetricInterceptor:    false,
		PathRelabel:                nil,
		cookieJar:                  nil,
		httpClient:                 nil,
	}, DefaultConfig()))
}

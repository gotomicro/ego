package egrpc

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/balancer/roundrobin"

	"github.com/gotomicro/ego/core/util/xtime"
)

func TestDefaultConfig(t *testing.T) {
	assert.True(t, reflect.DeepEqual(&Config{
		Addr:                         "",
		BalancerName:                 roundrobin.Name,
		OnFail:                       "panic",
		DialTimeout:                  time.Second * 3,
		ReadTimeout:                  xtime.Duration("1s"),
		SlowLogThreshold:             xtime.Duration("600ms"),
		EnableBlock:                  true,
		EnableOfficialGrpcLog:        false,
		EnableWithInsecure:           true,
		EnableMetricInterceptor:      true,
		EnableTraceInterceptor:       true,
		EnableAppNameInterceptor:     true,
		EnableTimeoutInterceptor:     true,
		EnableAccessInterceptor:      false,
		EnableAccessInterceptorReq:   false,
		EnableAccessInterceptorRes:   false,
		EnableFailOnNonTempDialError: true,
		EnableServiceConfig:          true,
		keepAlive:                    nil,
		dialOptions:                  nil,
		MaxCallRecvMsgSize:           DefaultMaxCallRecvMsgSize,
	}, DefaultConfig()))
}

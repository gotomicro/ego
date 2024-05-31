package ehttp

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	c := &Container{
		config: &Config{},
	}
	expectedAddr := "127.0.0.1:8080"
	expectedReadTimeOut := time.Duration(5)
	expectedSlowLogThreshold := time.Duration(5)
	expectedIdleConnTimeOut := time.Duration(5)

	WithAddr(expectedAddr)(c)
	WithDebug(true)(c)
	WithRawDebug(false)(c)
	WithReadTimeout(expectedReadTimeOut)(c)
	WithSlowLogThreshold(expectedSlowLogThreshold)(c)
	WithIdleConnTimeout(expectedIdleConnTimeOut)(c)
	WithMaxIdleConns(3)(c)
	WithMaxIdleConnsPerHost(3)(c)
	WithEnableTraceInterceptor(true)(c)
	WithEnableKeepAlives(true)(c)
	WithEnableMetricInterceptor(true)(c)
	WithEnableAccessInterceptor(true)(c)
	WithEnableAccessInterceptorRes(true)(c)
	WithPathRelabel("hello", "test")(c)
	WithJar(nil)(c)
	WithHTTPClient(nil)(c)

	assert.Equal(t, expectedAddr, c.config.Addr)
	assert.Equal(t, true, c.config.Debug)
	assert.Equal(t, false, c.config.RawDebug)
	assert.Equal(t, expectedReadTimeOut, c.config.ReadTimeout)
	assert.Equal(t, expectedSlowLogThreshold, c.config.SlowLogThreshold)
	assert.Equal(t, expectedIdleConnTimeOut, c.config.IdleConnTimeout)
	assert.Equal(t, 3, c.config.MaxIdleConns)
	assert.Equal(t, 3, c.config.MaxIdleConnsPerHost)
	assert.Equal(t, true, c.config.EnableTraceInterceptor)
	assert.Equal(t, true, c.config.EnableKeepAlives)
	assert.Equal(t, true, c.config.EnableMetricInterceptor)
	assert.Equal(t, true, c.config.EnableAccessInterceptor)
	assert.Equal(t, true, c.config.EnableAccessInterceptorRes)
	reflect.DeepEqual(Relabel{Match: "hello", Replacement: "test"}, c.config.PathRelabel)
	assert.Equal(t, nil, c.config.cookieJar)
	reflect.DeepEqual(nil, c.config.httpClient)
}

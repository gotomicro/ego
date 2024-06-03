package ehttp

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/client/ehttp/resolver"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/util/xtime"
)

func TestNewComponent(t *testing.T) {
	var logger = elog.DefaultLogger
	config := &Config{
		Addr:             "https://hello.com/",
		Debug:            false,
		RawDebug:         false,
		httpClient:       nil,
		ReadTimeout:      xtime.Duration("2s"),
		SlowLogThreshold: xtime.Duration("500ms"),
	}
	out := newComponent("test", config, logger)

	target, err := parseTarget(config.Addr)
	assert.NoError(t, err)
	assert.Equal(t, "https", target.Scheme)
	assert.Equal(t, "http", target.Protocol)
	assert.Equal(t, "", target.Endpoint)
	assert.Equal(t, "hello.com", target.Authority)

	config.httpClient = &http.Client{Transport: createTransport(config), Jar: config.cookieJar}
	cli := resty.NewWithClient(config.httpClient).
		SetDebug(config.RawDebug).
		SetTimeout(config.ReadTimeout).
		SetHeader("app", eapp.Name()).
		SetBaseURL(config.Addr)
	in := &Component{
		name:    "test",
		config:  config,
		logger:  logger,
		Client:  cli,
		builder: resolver.Get("https"),
	}
	assert.Equal(t, in.builder, out.builder)
	// assert.Equal(t, in.Client, out.Client)
}

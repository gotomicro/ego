package ehttp

import (
	"context"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/elog"
)

func TestLogAccess(t *testing.T) {
	name := "test"
	config := &Config{}
	logger := &elog.Component{}
	u, err := url.Parse("https://hello.com/xxx")
	assert.NoError(t, err)
	ctx := context.WithValue(context.Background(), urlKey{}, u)
	req := resty.New().R().SetContext(ctx)
	res := &resty.Response{}
	logAccess(name, config, logger, req, res, err)
	assert.NoError(t, err)
}

func TestBeg(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	ctx = context.WithValue(ctx, begKey{}, now)

	result := beg(ctx)
	assert.Equal(t, now, result)
}

type CustomResolver struct {
	Address string
}

func (r *CustomResolver) GetAddr() string {
	return r.Address
}

func TestFixedInterceptor(t *testing.T) {
	name := "test"
	config := &Config{}
	logger := &elog.Component{}
	builder := &CustomResolver{Address: "https://test.com"}

	client := resty.New()
	request := client.R()
	request.SetContext(context.Background())
	request.URL = "https://hello.com/world"
	middleware, _, _ := fixedInterceptor(name, config, logger, builder)

	// case 1
	config.Addr = ""
	err := middleware(client, request)
	assert.NoError(t, err)

	// case 2
	config.Addr = "https://xxxxx.com/xxx"
	err = middleware(client, request)
	assert.NoError(t, err)
	assert.Equal(t, "https://test.com", client.HostURL)
}

func TestFileWithLineNum(t *testing.T) {
	file := "/usr/local/go/src/testing/testing.go"
	got := fileWithLineNum()
	assert.True(t, true, strings.HasPrefix(got, file))
}

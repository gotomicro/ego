package ehttp

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/net/publicsuffix"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/elog"
)

// PackageName 设置包名
const PackageName = "client.ehttp"

// Component 组件
type Component struct {
	name   string
	config *Config
	logger *elog.Component
	*resty.Client
}

func newComponent(name string, config *Config, logger *elog.Component) *Component {
	// resty的默认方法，无法设置长连接个数，和是否开启长连接，这里重新构造http client。
	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	interceptors := []interceptor{fixedInterceptor, logInterceptor, metricInterceptor, traceInterceptor}
	cli := resty.NewWithClient(&http.Client{Transport: createTransport(config), Jar: cookieJar}).
		SetDebug(config.RawDebug).
		SetTimeout(config.ReadTimeout).
		SetHeader("app", eapp.Name()).
		SetBaseURL(config.Addr)
	for _, interceptor := range interceptors {
		onBefore, onAfter, onErr := interceptor(name, config, logger)
		if onBefore != nil {
			cli.OnBeforeRequest(onBefore)
		}
		if onAfter != nil {
			cli.OnAfterResponse(onAfter)
		}
		if onErr != nil {
			cli.OnError(onErr)
		}
	}

	return &Component{
		name:   name,
		config: config,
		logger: logger,
		Client: cli,
	}
}

func createTransport(config *Config) *http.Transport {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          config.MaxIdleConns,
		IdleConnTimeout:       config.IdleConnTimeout,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     !config.EnableKeepAlives,
		MaxIdleConnsPerHost:   config.MaxIdleConnsPerHost,
	}
}

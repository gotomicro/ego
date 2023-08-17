package ehttp

import (
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gotomicro/ego/client/ehttp/resolver"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/eregistry"
)

// PackageName 设置包名
const PackageName = "client.ehttp"

// Component 组件
type Component struct {
	name   string
	config *Config
	logger *elog.Component
	*resty.Client
	builder resolver.Builder
}

func newComponent(name string, config *Config, logger *elog.Component) *Component {
	// addr可以为空
	// 以下方法是为了支持k8s解析域名， k8s:///svc-user:9002 http接口，一定要是三斜线，跟gRPC统一
	egoTarget, err := parseTarget(config.Addr)
	if err != nil {
		elog.Panic("parse addr error", elog.FieldErr(err), elog.FieldKey(config.Addr))
	}
	// 这里的目的是为了，将k8s:// 替换为 http://，所以需要判断下是否为非HTTP，HTTPS。
	// 因为resty默认只要http和https的协议
	addr := config.Addr
	if egoTarget.Scheme != "http" && egoTarget.Scheme != "https" {
		// 因为内部协议，都是内网，所以直接替换为HTTP
		addr = strings.ReplaceAll(config.Addr, egoTarget.Scheme+"://", "http://")
	}
	builder := resolver.Get(egoTarget.Scheme)
	resolverBuild, err := builder.Build(addr)
	if err != nil {
		elog.Panic("build resolver error", elog.FieldErr(err), elog.FieldKey(config.Addr))
	}

	// resty的默认方法，无法设置长连接个数，和是否开启长连接，这里重新构造http client。
	interceptors := []interceptor{fixedInterceptor, logInterceptor, metricInterceptor, traceInterceptor}
	cli := resty.NewWithClient(&http.Client{Transport: createTransport(config), Jar: config.cookieJar}).
		SetDebug(config.RawDebug).
		SetTimeout(config.ReadTimeout).
		SetHeader("app", eapp.Name()).
		SetBaseURL(addr)
	for _, interceptor := range interceptors {
		onBefore, onAfter, onErr := interceptor(name, config, logger, resolverBuild)
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
		name:    name,
		config:  config,
		logger:  logger,
		Client:  cli,
		builder: builder,
	}
}

func parseTarget(addr string) (eregistry.Target, error) {
	target, err := url.Parse(addr)
	if err != nil {
		return eregistry.Target{}, err
	}
	endpoint := target.Path
	if endpoint == "" {
		endpoint = target.Opaque
	}
	endpoint = strings.TrimPrefix(endpoint, "/")

	egoTarget := eregistry.Target{
		Protocol:  eregistry.ProtocolHTTP,
		Scheme:    target.Scheme,
		Endpoint:  endpoint,
		Authority: target.Host,
	}
	return egoTarget, nil
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

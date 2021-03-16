package ehttp

import (
	"golang.org/x/net/publicsuffix"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/util/xdebug"
)

const PackageName = "client.ehttp"

type Component struct {
	name   string
	config *Config
	logger *elog.Component
	*resty.Client
}

func newComponent(name string, config *Config, logger *elog.Component) *Component {
	var logAccess = func(request *resty.Request, response *resty.Response, err error) {
		rr := request.RawRequest
		fullMethod := request.Method + "." + rr.URL.RequestURI() // GET./hello
		var (
			cost     time.Duration
			respBody string
		)
		if response != nil {
			cost = response.Time()
			respBody = string(response.Body())
		}
		if eapp.IsDevelopmentMode() {
			if err != nil {
				log.Println("http.response", xdebug.MakeReqResError(name, config.Addr, cost, fullMethod, err.Error()))
			} else {
				log.Println("http.response", xdebug.MakeReqResInfo(name, config.Addr, cost, fullMethod, respBody))
			}
		}

		var isSlowLog, isErrLog bool
		var fields = make([]elog.Field, 0, 15)

		fields = append(fields,
			elog.FieldMethod(fullMethod),
			elog.FieldName(name),
			elog.FieldCost(cost),
			elog.FieldAddr(rr.URL.Host),
		)

		if config.EnableAccessInterceptorRes {
			fields = append(fields, elog.FieldValueAny(respBody))
		}

		if err != nil {
			elog.FieldErr(err)
			if response == nil {
				// 无 response 的是连接超时等系统级错误
				fields = append(fields, elog.FieldEvent("error"))
				logger.Error("access", fields...)
			} else {
				logger.Warn("access", fields...)
			}
			isErrLog = true
		}

		if config.SlowLogThreshold > time.Duration(0) && cost > config.SlowLogThreshold {
			fields = append(fields, elog.FieldEvent("slow"))
			logger.Warn("access", fields...)
			isSlowLog = true
		}

		if config.EnableAccessInterceptor && !isSlowLog && !isErrLog {
			fields = append(fields, elog.FieldEvent("normal"))
			logger.Info("access", fields...)
		}
	}

	// resty的默认方法，无法设置长连接个数，和是否开启长连接，这里重新构造http client。
	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	restyClient := resty.NewWithClient(&http.Client{
		Transport: createTransport(config),
		Jar:       cookieJar,
	}).
		SetDebug(config.RawDebug).
		SetTimeout(config.ReadTimeout).
		SetHeader("app", eapp.Name()).
		OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
			logAccess(response.Request, response, nil)
			return nil
		}).
		OnError(func(req *resty.Request, err error) {
			if v, ok := err.(*resty.ResponseError); ok {
				logAccess(req, v.Response, v.Err)
			} else {
				logAccess(req, nil, err)
			}
		}).
		SetHostURL(config.Addr)

	return &Component{
		name:   name,
		config: config,
		logger: logger,
		Client: restyClient,
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

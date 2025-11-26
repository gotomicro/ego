package ehttp

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gotomicro/ego/core/transport"
	"github.com/spf13/cast"

	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/gotomicro/ego/client/ehttp/resolver"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/emetric"
	"github.com/gotomicro/ego/core/etrace"
	"github.com/gotomicro/ego/core/util/xdebug"
)

type interceptor func(name string, cfg *Config, logger *elog.Component, builder resolver.Resolver) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook)

func logAccess(name string, config *Config, logger *elog.Component, req *resty.Request, res *resty.Response, err error) {
	u := req.Context().Value(urlKey{}).(*url.URL)
	fullMethod := req.Method + "." + u.RequestURI() // GET./hello
	var cost = time.Since(beg(req.Context()))
	var respBody string
	if res != nil {
		respBody = string(res.Body())
	}
	if eapp.IsDevelopmentMode() {
		if err != nil {
			log.Println("http.response", xdebug.MakeReqAndResError(fileWithLineNum(), name, config.Addr, cost, fullMethod, err.Error()))
		} else {
			log.Println("http.response", xdebug.MakeReqAndResInfo(fileWithLineNum(), name, config.Addr, cost, fullMethod, respBody))
		}
	}

	loggerKeys := transport.CustomContextKeys()

	var fields = make([]elog.Field, 0, 16+transport.CustomContextKeysLength())
	fields = append(fields,
		elog.FieldMethod(fullMethod),
		elog.FieldName(name),
		elog.FieldCost(cost),
		elog.FieldAddr(u.Host),
	)

	event := "normal"

	// 支持自定义log
	for _, key := range loggerKeys {
		if value := req.Context().Value(key); value != nil {
			fields = append(fields, elog.FieldCustomKeyValue(key, cast.ToString(value)))
		}
	}

	// 开启了链路，那么就记录链路id
	if etrace.IsGlobalTracerRegistered() {
		fields = append(fields, elog.FieldTid(etrace.ExtractTraceID(req.Context())))
	}
	if config.EnableAccessInterceptor {
		if config.EnableAccessInterceptorReq {
			fields = append(fields, elog.Any("req", map[string]any{
				"metadata":      req.Header,
				"payload":       req.Body,
				"queryParam":    req.QueryParam,
				"formData":      req.FormData,
				"pathParams":    req.PathParams,
				"rawPathParams": req.RawPathParams,
			}))
		}

		if config.EnableAccessInterceptorRes {
			fields = append(fields, elog.Any("res", map[string]any{
				"metadata": res.Header(),
				"payload":  respBody,
			}))
		}
	}

	isSlowLog := false
	if config.SlowLogThreshold > time.Duration(0) && cost > config.SlowLogThreshold {
		event = "slow"
		isSlowLog = true
	}

	if err != nil {
		fields = append(fields, elog.FieldEvent(event), elog.FieldErr(err))
		if res == nil {
			// 无 res 的是连接超时等系统级错误
			logger.Error("access", fields...)
			return
		}
		logger.Warn("access", fields...)
		return
	}

	if config.EnableAccessInterceptor || isSlowLog {
		fields = append(fields, elog.FieldEvent(event))
		if isSlowLog {
			logger.Warn("access", fields...)
		} else {
			logger.Info("access", fields...)
		}
	}
}

// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
// https://blog.golang.org/context#TOC_3.2.
// https://golang.org/pkg/context/#WithValue ，这边文章说明了用struct，可以避免分配
type begKey struct{}
type urlKey struct{}

func beg(ctx context.Context) time.Time {
	begTime, _ := ctx.Value(begKey{}).(time.Time)
	return begTime
}

func fixedInterceptor(name string, config *Config, logger *elog.Component, builder resolver.Resolver) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook) {
	return func(cli *resty.Client, req *resty.Request) error {
		// 这个URL可能不准，每次请求都需要重复url.Parse()，会增加一定的性能损耗
		var concatURL string
		if config.Addr == "" {
			// 没有配置addr，host可能在url里面 (request.Get("http://xxx.com/xxx"))
			concatURL = req.URL
		} else {
			concatURL = strings.TrimRight(config.Addr, "/") + "/" + strings.TrimLeft(req.URL, "/")
		}
		u, err := url.Parse(concatURL)
		if err != nil {
			logger.Warn("invalid url", elog.String("concatURL", concatURL), elog.FieldErr(err))
			req.SetContext(context.WithValue(context.WithValue(req.Context(), begKey{}, time.Now()), urlKey{}, &url.URL{}))
			return err
		}
		if len(config.PathRelabel) > 0 {
			for _, relabel := range config.PathRelabel {
				if relabel.matchReg.MatchString(u.Path) {
					u.Path = relabel.Replacement
					break
				}
			}
		}
		// 只有存在，才会更新
		if builder.GetAddr() != "" {
			cli.HostURL = builder.GetAddr()
		}
		req.SetContext(context.WithValue(context.WithValue(req.Context(), begKey{}, time.Now()), urlKey{}, u))
		return nil
	}, nil, nil
}

func logInterceptor(name string, config *Config, logger *elog.Component, builder resolver.Resolver) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook) {
	loggerKeys := transport.CustomContextKeys()
	beforeFn := func(cli *resty.Client, req *resty.Request) error {
		// 增加header
		for _, key := range loggerKeys {
			if value := req.Context().Value(key); value != nil {
				req.SetHeader(key, cast.ToString(value))
			}
		}
		return nil
	}

	afterFn := func(cli *resty.Client, response *resty.Response) error {
		logAccess(name, config, logger, response.Request, response, nil)
		return nil
	}
	errorFn := func(req *resty.Request, err error) {
		if v, ok := err.(*resty.ResponseError); ok {
			logAccess(name, config, logger, req, v.Response, v.Err)
		} else {
			logAccess(name, config, logger, req, nil, err)
		}
	}
	return beforeFn, afterFn, errorFn
}

func metricInterceptor(name string, config *Config, logger *elog.Component, builder resolver.Resolver) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook) {
	if !config.EnableMetricInterceptor {
		return nil, nil, nil
	}
	addr := strings.TrimRight(config.Addr, "/")
	afterFn := func(cli *resty.Client, res *resty.Response) error {
		method := res.Request.Method + "." + res.Request.Context().Value(urlKey{}).(*url.URL).Path
		emetric.ClientHandleCounter.Inc(emetric.TypeHTTP, name, method, addr, http.StatusText(res.StatusCode()))
		emetric.ClientHandleHistogram.Observe(res.Time().Seconds(), emetric.TypeHTTP, name, method, addr)
		return nil
	}
	errorFn := func(req *resty.Request, err error) {
		method := req.Method + "." + req.Context().Value(urlKey{}).(*url.URL).Path
		if v, ok := err.(*resty.ResponseError); ok {
			emetric.ClientHandleCounter.Inc(emetric.TypeHTTP, name, method, addr, http.StatusText(v.Response.StatusCode()))
		} else {
			emetric.ClientHandleCounter.Inc(emetric.TypeHTTP, name, method, addr, "biz error")
		}
		emetric.ClientHandleHistogram.Observe(time.Since(beg(req.Context())).Seconds(), emetric.TypeHTTP, name, method, addr)
	}
	return nil, afterFn, errorFn
}

func traceInterceptor(name string, config *Config, logger *elog.Component, builder resolver.Resolver) (resty.RequestMiddleware, resty.ResponseMiddleware, resty.ErrorHook) {
	tracer := etrace.NewTracer(trace.SpanKindClient)
	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("http"),
	}
	beforeFn := func(cli *resty.Client, req *resty.Request) error {
		// 需要拿到header，才能将链路穿起来
		carrier := propagation.HeaderCarrier(req.Header)
		ctx, span := tracer.Start(req.Context(), req.Method, carrier, trace.WithAttributes(attrs...))
		span.SetAttributes(
			semconv.PeerServiceKey.String(name),
			semconv.HTTPMethodKey.String(req.Method),
			semconv.HTTPURLKey.String(req.URL),
		)
		req.SetContext(ctx)
		return nil
	}
	afterFn := func(cli *resty.Client, res *resty.Response) error {
		span := trace.SpanFromContext(res.Request.Context())
		span.SetAttributes(
			semconv.HTTPStatusCodeKey.Int64(int64(res.StatusCode())),
		)
		span.End()
		return nil
	}
	errorFn := func(req *resty.Request, err error) {
		span := trace.SpanFromContext(req.Context())
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}
	return beforeFn, afterFn, errorFn
}

func fileWithLineNum() string {
	// the second caller usually from internal, so set i start from 2
	for i := 2; i < 20; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if (!(strings.Contains(file, "ego") && strings.HasSuffix(file, "client/ehttp/interceptor.go")) && !strings.Contains(file, "go-resty/resty")) || strings.HasSuffix(file, "_test.go") {
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}
	return ""
}

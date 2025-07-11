package egin

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"github.com/google/cel-go/common/types"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	rpcpb "google.golang.org/genproto/googleapis/rpc/context/attribute_context"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/emetric"
	"github.com/gotomicro/ego/core/esentinel"
	"github.com/gotomicro/ego/core/etrace"
	"github.com/gotomicro/ego/core/transport"
	"github.com/gotomicro/ego/internal/tools"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// extractAPP 提取header头中的app信息
func extractAPP(ctx *gin.Context) string {
	return ctx.Request.Header.Get("app")
}

type resWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (g *resWriter) Write(data []byte) (int, error) {
	n, e := g.body.Write(data)
	if e != nil {
		return n, e
	}
	return g.ResponseWriter.Write(data)
}

func (g *resWriter) WriteString(s string) (int, error) {
	n, e := g.body.WriteString(s)
	if e != nil {
		return n, e
	}
	return g.ResponseWriter.WriteString(s)
}

func copyHeaders(headers http.Header) http.Header {
	nh := http.Header{}
	for k, v := range headers {
		nh[k] = v
	}
	return nh
}

// timeout middleware wraps the request context with a timeout
func timeoutMiddleware(timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		// 若无自定义超时设置，默认设置超时
		_, ok := c.Request.Context().Deadline()
		if ok {
			c.Next()
			return
		}

		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer func() {
			// check if context timeout was reached
			if ctx.Err() == context.DeadlineExceeded {

				// write response and abort the request
				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				c.Abort()
			}

			// cancel to clear resources after finished
			cancel()
		}()

		// replace request with context wrapped request
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// defaultServerInterceptor 默认拦截器，包含日志记录、Recover、监控功能
// 监控放里面是因为，例如panic会改写http status。这样才能统计准确
func (c *Container) defaultServerInterceptor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var beg = time.Now()
		var rw *resWriter
		var rb bytes.Buffer

		// 只有开启了EnableAccessInterceptorRes时拷贝request body
		// 也可以直接使用econf.Sub(c.name).GetBool("EnableAccessInterceptorReq")，不过从econf动态查找配置性能可能会比较差，暂时先用锁代替
		c.config.mu.RLock()
		if c.config.EnableAccessInterceptorReq || c.config.AccessInterceptorReqResFilter != "" {
			ctx.Request.Body = io.NopCloser(io.TeeReader(ctx.Request.Body, &rb))
		}
		// 只有开启了EnableAccessInterceptorRes时才替换response writer
		if c.config.EnableAccessInterceptorRes || c.config.AccessInterceptorReqResFilter != "" {
			rw = &resWriter{ctx.Writer, &bytes.Buffer{}}
			ctx.Writer = rw
		}
		c.config.mu.RUnlock()

		// 为了性能考虑，如果要加日志字段，需要改变slice大小
		loggerKeys := transport.CustomContextKeys()
		var fields = make([]elog.Field, 0, 20+len(loggerKeys))
		var brokenPipe bool
		var event = "normal"

		// 必须在defer外层，因为要赋值，替换ctx
		// 只有在环境变量里的自定义header，才会写入到context value里
		for _, key := range loggerKeys {
			// 赋值context
			getHeaderValue(ctx, key, c.config.EnableTrustedCustomHeader)
		}

		defer func() {
			cost := time.Since(beg)
			fields = append(fields,
				elog.FieldKey(ctx.Request.Method), // GET, POST
				elog.FieldCost(cost),
				elog.FieldMethod(ctx.Request.Method+"."+ctx.FullPath()),
				elog.FieldAddr(ctx.Request.URL.RequestURI()),
				elog.FieldIP(ctx.ClientIP()),
				elog.FieldSize(int32(ctx.Writer.Size())),
				elog.FieldPeerIP(getPeerIP(ctx.Request.RemoteAddr)),
				elog.FieldPeerName(getPeerName(ctx)),
			)

			for _, key := range loggerKeys {
				if value := tools.ContextValue(ctx.Request.Context(), key); value != "" {
					fields = append(fields, elog.FieldCustomKeyValue(key, value))
				}
			}

			for _, key := range transport.CustomHeaderKeys() {
				if value := tools.ContextValue(ctx.Request.Context(), key); value != "" {
					// x-expose 需要在这里获取
					if strings.HasPrefix(key, eapp.EgoHeaderExpose()) {
						// 设置到ctx response header
						ctx.Writer.Header().Set(key, value)
					}
				}
			}

			if etrace.IsGlobalTracerRegistered() {
				fields = append(fields, elog.FieldTid(etrace.ExtractTraceID(ctx.Request.Context())))
			}

			c.config.mu.RLock()
			if c.config.EnableAccessInterceptorReq || c.config.EnableAccessInterceptorRes {
				out := c.checkFilter(ctx.Request, rw)

				if c.config.EnableAccessInterceptorReq && out {
					if len(rb.String()) > c.config.AccessInterceptorReqMaxLength {
						fields = append(fields, elog.Any("req", map[string]interface{}{
							"metadata": copyHeaders(ctx.Request.Header),
							"payload":  rb.String()[:c.config.AccessInterceptorReqMaxLength] + "...",
						}))
					} else {
						fields = append(fields, elog.Any("req", map[string]interface{}{
							"metadata": copyHeaders(ctx.Request.Header),
							"payload":  rb.String(),
						}))
					}
				}
				if c.config.EnableAccessInterceptorRes && out {
					if len(rw.body.String()) > c.config.AccessInterceptorResMaxLength {
						fields = append(fields, elog.Any("res", map[string]interface{}{
							"metadata": copyHeaders(ctx.Request.Header),
							"payload":  rw.body.String()[:c.config.AccessInterceptorResMaxLength] + "...",
						}))
					} else {
						fields = append(fields, elog.Any("res", map[string]interface{}{
							"metadata": copyHeaders(ctx.Writer.Header()),
							"payload":  rw.body.String(),
						}))
					}
				}
			}
			c.config.mu.RUnlock()

			// slow log
			isSlowLog := false
			if c.config.SlowLogThreshold > time.Duration(0) && c.config.SlowLogThreshold < cost {
				// 非长连接模式下，记入warn慢日志
				if ctx.GetHeader("Accept") != "text/event-stream" {
					isSlowLog = true
					event = "slow"
				}
			}

			if rec := recover(); rec != nil {
				if ne, ok := rec.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// BrokenPipe 使用用户的status
				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					ctx.Error(rec.(error)) // nolint: errcheck
					ctx.Abort()
				} else {
					// 如果不是，默认使用500错误码
					if c.config.recoveryFunc == nil {
						c.config.recoveryFunc = defaultRecoveryFunc
					}
					c.config.recoveryFunc(ctx, rec)
				}

				stackInfo := stack(3)
				fields = append(fields,
					elog.FieldEvent(event),
					elog.FieldType("recover"),
					zap.ByteString("stack", stackInfo),
					elog.FieldErrAny(rec),
					elog.FieldCode(int32(ctx.Writer.Status())),
					elog.FieldUniformCode(int32(ctx.Writer.Status())),
				)
				c.metricServerInterceptor(ctx, cost)
				// broken pipe 是warning
				if brokenPipe {
					c.logger.Warn("access", fields...)
				} else {
					c.logger.Error("access", fields...)
				}
				return
			}
			// todo 如果不记录日志的时候，应该早点return
			if c.config.EnableAccessInterceptor || isSlowLog {
				fields = append(fields,
					elog.FieldEvent(event),
					elog.FieldCode(int32(ctx.Writer.Status())),
					elog.FieldUniformCode(int32(ctx.Writer.Status())),
				)
				if errStr := ctx.Errors.ByType(gin.ErrorTypePrivate).String(); errStr != "" {
					fields = append(fields, elog.FieldErrAny(errStr))
				}
				if isSlowLog {
					c.logger.Warn("access", fields...)
				} else {
					c.logger.Info("access", fields...)
				}
			}
			c.metricServerInterceptor(ctx, cost)
		}()
		ctx.Next()
	}
}

// func copyBody(r io.Reader, w io.Writer) io.ReadCloser {
//	return os.NopCloser(io.TeeReader(r, w))
// }

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func (c *Container) metricServerInterceptor(ctx *gin.Context, cost time.Duration) {
	if !c.config.EnableMetricInterceptor {
		return
	}

	host := ctx.Request.Host
	method := ctx.Request.Method + "." + ctx.FullPath()
	app := extractAPP(ctx)
	emetric.ServerStartedCounter.Inc(emetric.TypeHTTP, method, app, host)
	// HandleHistogram的单位是s，需要用s单位
	emetric.ServerHandleHistogram.ObserveWithExemplar(cost.Seconds(), prometheus.Labels{
		"tid": etrace.ExtractTraceID(ctx.Request.Context()),
	}, emetric.TypeHTTP, method, app, host)
	emetric.ServerHandleCounter.Inc(emetric.TypeHTTP, method, app, http.StatusText(ctx.Writer.Status()), strconv.Itoa(ctx.Writer.Status()), host)
}

// todo 如果业务崩了，logger recover
func traceServerInterceptor(compatibleTraceFunc ...func(http.Header)) gin.HandlerFunc {
	tracer := etrace.NewTracer(trace.SpanKindServer)
	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("http"),
	}
	return func(c *gin.Context) {
		// 执行自定义的trace处理函数
		if len(compatibleTraceFunc) > 0 && compatibleTraceFunc[0] != nil {
			compatibleTraceFunc[0](c.Request.Header)
		}

		// 该方法会在v0.9.0移除
		// etrace.CompatibleExtractHTTPTraceID(c.Request.Header)
		ctx, span := tracer.Start(c.Request.Context(), c.Request.Method+"."+c.FullPath(), propagation.HeaderCarrier(c.Request.Header), trace.WithAttributes(attrs...))
		span.SetAttributes(
			semconv.HTTPURLKey.String(c.Request.URL.String()),
			semconv.HTTPTargetKey.String(c.Request.URL.Path),
			semconv.HTTPMethodKey.String(c.Request.Method),
			semconv.HTTPUserAgentKey.String(c.Request.UserAgent()),
			semconv.HTTPClientIPKey.String(c.ClientIP()),
			etrace.CustomTag("http.full_path", c.FullPath()),
		)
		c.Request = c.Request.WithContext(ctx)
		c.Header(eapp.EgoTraceIDName(), span.SpanContext().TraceID().String())
		c.Next()
		span.SetAttributes(
			semconv.HTTPStatusCodeKey.Int64(int64(c.Writer.Status())),
		)
		span.End()
	}
}

// sentinelMiddleware returns new gin.HandlerFunc
// Default resource name is {method}:{path}, such as "GET:/api/users/:id"
// Default block fallback is returning 429 code
// Define your own behavior by setting options
func (c *Container) sentinelMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resourceName := ctx.Request.Method + "." + ctx.FullPath()

		if c.config.resourceExtract != nil {
			resourceName = c.config.resourceExtract(ctx)
		}

		if !esentinel.IsResExist(resourceName) {
			ctx.Next()
			return
		}

		entry, err := sentinel.Entry(
			resourceName,
			sentinel.WithResourceType(base.ResTypeWeb),
			sentinel.WithTrafficType(base.Inbound),
		)

		if err != nil {
			if c.config.blockFallback != nil {
				c.config.blockFallback(ctx)
			} else {
				ctx.AbortWithStatus(http.StatusTooManyRequests)
			}
			return
		}

		defer entry.Exit()

		ctx.Next()
	}
}

func getPeerIP(addr string) string {
	addSlice := strings.Split(addr, ":")
	if len(addSlice) > 1 {
		return addSlice[0]
	}
	return ""
}

func getPeerName(c *gin.Context) string {
	value := c.GetHeader("app")
	return value
}

func getHeaderValue(c *gin.Context, key string, enableTrustedCustomHeader bool) string {
	if key == "" {
		return ""
	}
	// 通常HTTP在外网，例如自定义Header： X-Ego-Uid 不可信任
	if !enableTrustedCustomHeader {
		return ""
	}
	value := c.GetHeader(key)
	if value != "" {
		// 如果信任该Header，将header数据赋值到context上
		c.Request = c.Request.WithContext(transport.WithValue(c.Request.Context(), key, value))
	}
	return value
}

func convert2googleResponse(rw *resWriter) *rpcpb.AttributeContext_Response {
	return &rpcpb.AttributeContext_Response{
		Code:    int64(rw.Status()),
		Headers: convertHeader(rw.Header()),
		Time:    timestamppb.New(time.Now()),
	}
}

func convert2googleRequest(r *http.Request) *rpcpb.AttributeContext_Request {
	return &rpcpb.AttributeContext_Request{
		Method:  r.Method,
		Headers: convertHeader(r.Header),
		Path:    r.URL.Path,
		Host:    r.Host,
		Scheme:  r.URL.Scheme,
		Query:   r.URL.RawQuery,
		Time:    timestamppb.New(time.Now()),
	}
}

func convertHeader(headers http.Header) map[string]string {
	h := make(map[string]string)
	for name, val := range headers {
		h[strings.ToLower(name)] = strings.Join(val, ";")
	}
	return h
}

func (c *Container) checkFilter(req *http.Request, rw *resWriter) bool {
	if c.config.aiReqResCelPrg == nil {
		return true
	}
	out, _, err := c.config.aiReqResCelPrg.Eval(map[string]interface{}{
		"request":  convert2googleRequest(req),
		"response": convert2googleResponse(rw),
	})
	if err != nil {
		c.logger.Warn("cel eval fail", elog.FieldErr(err))
	}
	return out == types.True
}

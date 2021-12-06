package egin

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/emetric"
	"github.com/gotomicro/ego/core/etrace"
	"github.com/gotomicro/ego/core/transport"
	"github.com/gotomicro/ego/internal/tools"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
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

// defaultServerInterceptor 默认拦截器，包含日志记录、Recover等功能
func (c *Container) defaultServerInterceptor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var beg = time.Now()
		var rw *resWriter
		var rb bytes.Buffer

		// 只有开启了EnableAccessInterceptorRes时拷贝request body
		if c.config.EnableAccessInterceptorReq {
			ctx.Request.Body = ioutil.NopCloser(io.TeeReader(ctx.Request.Body, &rb))
		}
		// 只有开启了EnableAccessInterceptorRes时才替换response writer
		if c.config.EnableAccessInterceptorRes {
			rw = &resWriter{ctx.Writer, &bytes.Buffer{}}
			ctx.Writer = rw
		}

		// 为了性能考虑，如果要加日志字段，需要改变slice大小
		loggerKeys := transport.CustomContextKeys()
		var fields = make([]elog.Field, 0, 20+len(loggerKeys))
		var brokenPipe bool
		var event = "normal"

		// 必须在defer外层，因为要赋值，替换ctx
		for _, key := range loggerKeys {
			// 赋值context
			getHeaderValue(ctx, key, c.config.EnableTrustedCustomHeader)
		}

		defer func() {
			cost := time.Since(beg)
			fields = append(fields,
				elog.FieldType("http"), // GET, POST
				elog.FieldCost(cost),
				elog.FieldMethod(ctx.Request.Method+"."+ctx.FullPath()),
				elog.FieldAddr(ctx.Request.URL.RequestURI()),
				elog.FieldIP(ctx.ClientIP()),
				elog.FieldSize(int32(ctx.Writer.Size())),
				elog.FieldPeerIP(getPeerIP(ctx.Request.RemoteAddr)),
			)

			for _, key := range loggerKeys {
				if value := tools.ContextValue(ctx.Request.Context(), key); value != "" {
					fields = append(fields, elog.FieldCustomKeyValue(key, value))
				}
			}

			if c.config.EnableTraceInterceptor && etrace.IsGlobalTracerRegistered() {
				fields = append(fields, elog.FieldTid(etrace.ExtractTraceID(ctx.Request.Context())))
			}

			if c.config.EnableAccessInterceptorReq {
				fields = append(fields, elog.Any("req", map[string]interface{}{
					"metadata": copyHeaders(ctx.Request.Header),
					"payload":  rb.String(),
				}))
			}

			if c.config.EnableAccessInterceptorRes {
				fields = append(fields, elog.Any("res", map[string]interface{}{
					"metadata": copyHeaders(ctx.Writer.Header()),
					"payload":  rw.body.String(),
				}))
			}

			// slow log
			if c.config.SlowLogThreshold > time.Duration(0) && c.config.SlowLogThreshold < cost {
				c.logger.Warn("slow", fields...)
			}

			if rec := recover(); rec != nil {
				if ne, ok := rec.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					ctx.Error(rec.(error)) // nolint: errcheck
					ctx.Abort()
				} else {
					ctx.AbortWithStatus(http.StatusInternalServerError)
				}

				event = "recover"
				stackInfo := stack(3)
				fields = append(fields,
					elog.FieldEvent(event),
					zap.ByteString("stack", stackInfo),
					elog.FieldErrAny(rec),
					elog.FieldCode(int32(ctx.Writer.Status())),
					elog.FieldUniformCode(int32(ctx.Writer.Status())),
				)
				c.logger.Error("access", fields...)
				return
			}
			// todo 如果不记录日志的时候，应该早点return
			if c.config.EnableAccessInterceptor {
				fields = append(fields,
					elog.FieldEvent(event),
					elog.FieldErrAny(ctx.Errors.ByType(gin.ErrorTypePrivate).String()),
					elog.FieldCode(int32(ctx.Writer.Status())),
				)
				c.logger.Info("access", fields...)
			}
		}()
		ctx.Next()
	}
}

//func copyBody(r io.Reader, w io.Writer) io.ReadCloser {
//	return ioutil.NopCloser(io.TeeReader(r, w))
//}

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
			data, err := ioutil.ReadFile(file)
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

func metricServerInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		beg := time.Now()
		c.Next()
		emetric.ServerHandleHistogram.Observe(time.Since(beg).Seconds(), emetric.TypeHTTP, c.Request.Method+"."+c.FullPath(), extractAPP(c))
		emetric.ServerHandleCounter.Inc(emetric.TypeHTTP, c.Request.Method+"."+c.FullPath(), extractAPP(c), http.StatusText(c.Writer.Status()), http.StatusText(c.Writer.Status()))
	}
}

func traceServerInterceptor() gin.HandlerFunc {
	tracer := etrace.NewTracer(trace.SpanKindServer)
	return func(c *gin.Context) {
		// 该方法会在v0.9.0移除
		etrace.CompatibleExtractHttpTraceId(c.Request.Header)
		ctx, span := tracer.Start(c.Request.Context(), c.Request.Method+"."+c.FullPath(), propagation.HeaderCarrier(c.Request.Header))
		span.SetAttributes(
			etrace.TagComponent("http"),
			etrace.CustomTag("http.url", c.Request.URL.Path),
			etrace.CustomTag("http.method", c.Request.Method),
			etrace.CustomTag("peer.ipv4", c.ClientIP()),
		)
		c.Request = c.Request.WithContext(ctx)
		defer span.End()
		c.Header(eapp.EgoTraceIDName(), span.SpanContext().TraceID().String())
		c.Next()
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

// 获取对端ip
func getPeerIP(addr string) string {
	addSlice := strings.Split(addr, ":")
	if len(addSlice) > 1 {
		return addSlice[0]
	}
	return ""
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

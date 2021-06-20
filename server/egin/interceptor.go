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

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/cast"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"

	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/emetric"
	"github.com/gotomicro/ego/core/etrace"
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
func defaultServerInterceptor(logger *elog.Component, config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var beg = time.Now()
		var rw *resWriter
		var rb bytes.Buffer

		// 只有开启了EnableAccessInterceptorRes时拷贝request body
		if config.EnableAccessInterceptorReq {
			c.Request.Body = ioutil.NopCloser(io.TeeReader(c.Request.Body, &rb))
		}
		// 只有开启了EnableAccessInterceptorRes时才替换response writer
		if config.EnableAccessInterceptorRes {
			rw = &resWriter{c.Writer, &bytes.Buffer{}}
			c.Writer = rw
		}

		// 为了性能考虑，如果要加日志字段，需要改变slice大小
		loggerKeys := eapp.EgoLoggerKeys()
		var fields = make([]elog.Field, 0, 20+len(loggerKeys))
		var brokenPipe bool
		var event = "normal"
		defer func() {
			cost := time.Since(beg)
			fields = append(fields,
				elog.FieldType("http"), // GET, POST
				elog.FieldCost(cost),
				elog.FieldMethod(c.Request.Method+"."+c.FullPath()),
				elog.FieldAddr(c.Request.URL.RequestURI()),
				elog.FieldIP(c.ClientIP()),
				elog.FieldSize(int32(c.Writer.Size())),
				elog.FieldPeerIP(getPeerIP(c.Request.RemoteAddr)),
			)

			if config.EnableTraceInterceptor && opentracing.IsGlobalTracerRegistered() {
				fields = append(fields, elog.FieldTid(etrace.ExtractTraceID(c.Request.Context())))
			}

			if config.EnableAccessInterceptorReq {
				fields = append(fields, elog.Any("req", map[string]interface{}{
					"metadata": copyHeaders(c.Request.Header),
					"body":     rb.String(),
				}))
			}

			if config.EnableAccessInterceptorRes {
				fields = append(fields, elog.Any("res", map[string]interface{}{
					"metadata": copyHeaders(c.Writer.Header()),
					"body":     rw.body.String(),
				}))
			}

			for _, v := range loggerKeys {
				if value := getContextValue(v, c); value != "" {
					fields = append(fields, elog.FieldCustomKeyValue(v, value))
				}
			}

			// slow log
			if config.SlowLogThreshold > time.Duration(0) && config.SlowLogThreshold < cost {
				logger.Warn("slow", fields...)
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
					c.Error(rec.(error)) // nolint: errcheck
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}

				event = "recover"
				stackInfo := stack(3)

				fields = append(fields,
					elog.FieldEvent(event),
					zap.ByteString("stack", stackInfo),
					elog.FieldErrAny(rec),
					elog.FieldCode(int32(c.Writer.Status())),
				)
				logger.Error("access", fields...)
				return
			}

			fields = append(fields,
				elog.FieldEvent(event),
				elog.FieldErrAny(c.Errors.ByType(gin.ErrorTypePrivate).String()),
				elog.FieldCode(int32(c.Writer.Status())),
			)
			logger.Info("access", fields...)
		}()
		c.Next()
	}
}

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
		emetric.ServerHandleCounter.Inc(emetric.TypeHTTP, c.Request.Method+"."+c.FullPath(), extractAPP(c), http.StatusText(c.Writer.Status()))
	}
}

func traceServerInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := etrace.StartSpanFromContext(
			c.Request.Context(),
			c.Request.Method+"."+c.FullPath(),
			etrace.TagComponent("http"),
			etrace.TagSpanKind("server"),
			etrace.HeaderExtractor(c.Request.Header),
			etrace.CustomTag("http.url", c.Request.URL.Path),
			etrace.CustomTag("http.method", c.Request.Method),
			etrace.CustomTag("peer.ipv4", c.ClientIP()),
		)
		c.Request = c.Request.WithContext(ctx)
		defer span.Finish()
		// 判断了全局jaeger的设置，所以这里一定能够断言为jaeger
		c.Header(eapp.EgoTraceIDName(), span.(*jaeger.Span).Context().(jaeger.SpanContext).TraceID().String())
		c.Next()
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

func getContextValue(key string, c *gin.Context) string {
	if key == "" {
		return ""
	}
	// 用Request.Context，因为这个是原生的HTTP，会往下传递链路，所以复用该Context传递的信息
	val := cast.ToString(c.Request.Context().Value(key))
	if val == "" {
		return c.GetHeader(key)
	}
	return val
}

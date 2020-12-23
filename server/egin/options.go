package egin

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/emetric"
	"github.com/gotomicro/ego/core/etrace"
	"go.uber.org/zap"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func extractAID(ctx *gin.Context) string {
	return ctx.Request.Header.Get("AID")
}

func recoverMiddleware(logger *elog.Component, slowLogThreshold time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var beg = time.Now()
		var fields = make([]elog.Field, 0, 8)
		var brokenPipe bool
		var event = "normal"
		defer func() {
			cost := time.Since(beg)

			// slow log
			if slowLogThreshold > time.Duration(0) && slowLogThreshold < cost {
				event = "slow"
			}

			fields = append(fields,
				elog.FieldCost(cost),
				elog.FieldType(c.Request.Method), // GET, POST
				elog.FieldMethod(c.Request.URL.Path),

				elog.FieldIp(c.ClientIP()),
				elog.FieldSize(int32(c.Writer.Size())),
				elog.FieldPeerIP(getPeerIP(c.Request.RemoteAddr)),
			)

			if rec := recover(); rec != nil {
				if ne, ok := rec.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				event = "recover"
				var err = rec.(error)
				fields = append(fields,
					elog.FieldEvent(event),
					zap.ByteString("stack", stack(3)),
					elog.FieldErr(err),
					elog.FieldCode(int32(http.StatusInternalServerError)),
				)
				logger.Error("access", fields...)
				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					c.Error(err) // nolint: errcheck
					c.Abort()
					return
				}
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			// httpRequest, _ := httputil.DumpRequest(c.Request, false)
			// fields = append(fields, zap.ByteString("request", httpRequest))
			fields = append(fields,
				elog.FieldEvent(event),
				zap.String("err", c.Errors.ByType(gin.ErrorTypePrivate).String()),
				elog.FieldCode(int32(c.Writer.Status())),
			)

			if event == "slow" {
				logger.Warn("access", fields...)
			} else {
				logger.Info("access", fields...)
			}
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
		emetric.ServerHandleHistogram.Observe(time.Since(beg).Seconds(), emetric.TypeHTTP, c.Request.Method+"."+c.Request.URL.Path, extractAID(c))
		emetric.ServerHandleCounter.Inc(emetric.TypeHTTP, c.Request.Method+"."+c.Request.URL.Path, extractAID(c), http.StatusText(c.Writer.Status()))
		return
	}
}

func traceServerInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := etrace.StartSpanFromContext(
			c.Request.Context(),
			c.Request.Method+" "+c.Request.URL.Path,
			etrace.TagComponent("http"),
			etrace.TagSpanKind("server"),
			etrace.HeaderExtractor(c.Request.Header),
			etrace.CustomTag("http.url", c.Request.URL.Path),
			etrace.CustomTag("http.method", c.Request.Method),
			etrace.CustomTag("peer.ipv4", c.ClientIP()),
		)
		c.Request = c.Request.WithContext(ctx)
		defer span.Finish()
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

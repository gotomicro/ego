package egin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"syscall"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/transport"

	"github.com/gotomicro/ego/core/elog"
)

func TestPanicInHandler(t *testing.T) {
	router := gin.New()
	// 使用非异步日志
	logger := elog.DefaultContainer().Build(
		elog.WithDebug(false),
		elog.WithEnableAddCaller(true),
		elog.WithEnableAsync(false),
	)
	container := DefaultContainer()
	container.Build(WithLogger(logger))

	// 使用recover组件
	router.Use(container.defaultServerInterceptor())
	router.GET("/recovery", func(c *gin.Context) {
		c.Status(200)
		panic("we have a panic")
	})
	// 调用触发panic的接口
	w := performRequest(router, "GET", "/recovery")
	logged, err := ioutil.ReadFile(path.Join(logger.ConfigDir(), logger.ConfigName()))
	fmt.Printf("logged--------------->"+"%+v\n", string(logged))
	assert.Nil(t, err)
	// 虽然程序里返回200，只要panic就会为500
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var m map[string]interface{}
	n := strings.Index(string(logged), "{")
	err = json.Unmarshal(logged[n:], &m)
	assert.NoError(t, err)
	assert.Contains(t, m["event"], `recover`)
	assert.Contains(t, string(logged), "we have a panic")
	assert.Contains(t, m["method"], `GET./recovery`)
	assert.Contains(t, string(logged), "500")
	assert.Contains(t, string(logged), t.Name())
	os.Remove(path.Join(logger.ConfigDir(), logger.ConfigName()))
}

func TestPanicInCustomHandler(t *testing.T) {
	router := gin.New()
	// 使用非异步日志
	logger := elog.DefaultContainer().Build(
		elog.WithDebug(false),
		elog.WithEnableAddCaller(true),
		elog.WithEnableAsync(false),
	)

	// 自定义 recover
	var recoverFunc gin.RecoveryFunc = func(ctx *gin.Context, err interface{}) {
		ctx.String(http.StatusInternalServerError, "%v", err)
		ctx.Abort()
	}

	container := DefaultContainer()
	container.Build(WithLogger(logger), WithRecoveryFunc(recoverFunc))

	// 使用recover组件
	panicMessage := "we have a panic"
	router.Use(container.defaultServerInterceptor())
	router.GET("/recovery", func(_ *gin.Context) {
		panic(panicMessage)
	})
	// 调用触发panic的接口
	w := performRequest(router, "GET", "/recovery")
	logged, err := ioutil.ReadFile(path.Join(logger.ConfigDir(), logger.ConfigName()))
	fmt.Printf("logged--------------->%+v\n", string(logged))
	assert.Nil(t, err)
	// TEST
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, panicMessage, w.Body.String())
	var m map[string]interface{}
	n := strings.Index(string(logged), "{")
	err = json.Unmarshal(logged[n:], &m)
	assert.NoError(t, err)
	assert.Equal(t, m["event"], `recover`)
	assert.Equal(t, m["error"], panicMessage)
	assert.Equal(t, m["method"], `GET./recovery`)
	assert.Contains(t, m["stack"], t.Name())
	os.Remove(path.Join(logger.ConfigDir(), logger.ConfigName()))
}

func TestPanicWithBrokenPipe(t *testing.T) {
	const expectCode = 204

	expectMsgs := map[syscall.Errno]string{
		syscall.EPIPE:      "broken pipe",
		syscall.ECONNRESET: "connection reset by peer",
	}

	for errno, expectMsg := range expectMsgs {
		t.Run(expectMsg, func(t *testing.T) {
			router := gin.New()
			logger := elog.DefaultContainer().Build(
				elog.WithDebug(false),
				elog.WithEnableAddCaller(true),
				elog.WithEnableAsync(false),
			)
			container := DefaultContainer()
			container.Build(WithLogger(logger))
			router.Use(container.defaultServerInterceptor())
			router.GET("/recovery", func(c *gin.Context) {
				// Start writing response
				c.Header("X-Test", "Value")
				c.Status(expectCode)

				// Oops. Client connection closed
				e := &net.OpError{Err: &os.SyscallError{Err: errno}}
				panic(e)
			})
			// RUN
			w := performRequest(router, "GET", "/recovery")
			logged, err := ioutil.ReadFile(path.Join(logger.ConfigDir(), logger.ConfigName()))
			assert.Nil(t, err)
			// TEST
			assert.Equal(t, expectCode, w.Code)
			assert.Contains(t, string(logged), `"code":204`)
			assert.Contains(t, string(logged), expectMsg)
			os.Remove(path.Join(logger.ConfigDir(), logger.ConfigName()))
		})
	}
}

func TestPrometheus(t *testing.T) {
	// 1 获取prometheus的handler的数据
	ts := httptest.NewServer(promhttp.Handler())
	defer ts.Close()
	container := DefaultContainer()
	cmp := container.Build()
	cmp.GET("/hello", func(_ *gin.Context) {})
	// RUN
	w := performRequest(cmp, "GET", "/hello")
	assert.Equal(t, 200, w.Code)
	pc := ts.Client()
	res, err := pc.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	text, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	err = res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, string(text), `ego_server_handle_seconds_count{method="GET./hello",peer="",rpc_service="example.com",type="http"}`)
	assert.Contains(t, string(text), `ego_server_started_total{method="GET./hello",peer="",rpc_service="example.com",type="http"}`)
}

type header struct {
	Key   string
	Value string
}

func performRequest(r http.Handler, method, path string, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func Test_getHeaderValue(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/chat", nil)
	c.Request.Header.Set("X-Ego-Uid", "9527")
	value := getHeaderValue(c, "X-Ego-Uid", true)
	assert.Equal(t, "9527", value)
}

func Test_getHeaderAssignValue(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/chat", nil)
	c.Request.Header.Set("X-Ego-Uid", "9527")
	value := getHeaderValue(c, "X-Ego-Uid", true)
	assert.Equal(t, "9527", value)

	value2 := transport.Value(c.Request.Context(), "X-Ego-Uid")
	assert.Equal(t, "9527", value2)
}

func Test_getPeerIp(t *testing.T) {
	addr := getPeerIP("192.168.1.1:50085")
	assert.Equal(t, "192.168.1.1", addr)
}

func Test_copyBody(t *testing.T) {
	src := []byte("hello, world")
	dst := make([]byte, len(src))
	rb := bytes.NewBuffer(src)
	var wb bytes.Buffer
	r := io.TeeReader(rb, &wb)
	if n, err := io.ReadFull(r, dst); err != nil || n != len(src) {
		t.Fatalf("ReadFull(r, dst) = %d, %v; want %d, nil", n, err, len(src))
	}
}

package egin

import (
	"bytes"
	"encoding/json"
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
	"github.com/gotomicro/ego/core/transport"
	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/elog"
)

func TestPanicInHandler(t *testing.T) {
	router := gin.New()
	logger := elog.DefaultContainer().Build(
		elog.WithDebug(false),
		elog.WithEnableAddCaller(true),
		elog.WithEnableAsync(false),
	)
	router.Use(defaultServerInterceptor(logger, DefaultConfig()))
	router.GET("/recovery", func(_ *gin.Context) {
		panic("we have a panic")
	})
	// RUN
	w := performRequest(router, "GET", "/recovery")
	logged, err := ioutil.ReadFile(path.Join(logger.ConfigDir(), logger.ConfigName()))
	assert.Nil(t, err)
	// TEST
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var m map[string]interface{}
	n := strings.Index(string(logged), "{")
	err = json.Unmarshal(logged[n:], &m)
	assert.NoError(t, err)
	assert.Contains(t, m["event"], `recover`)
	assert.Contains(t, string(logged), "we have a panic")
	assert.Contains(t, m["method"], `GET./recovery`)
	assert.Contains(t, string(logged), t.Name())
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
			router.Use(defaultServerInterceptor(logger, DefaultConfig()))
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

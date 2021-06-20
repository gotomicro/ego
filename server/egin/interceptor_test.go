package egin

import (
	"encoding/json"
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
	logged, err := ioutil.ReadFile(path.Join(logger.GetConfigDir(), logger.GetConfigName()))
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
	os.Remove(path.Join(logger.GetConfigDir(), logger.GetConfigName()))
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
			logged, err := ioutil.ReadFile(path.Join(logger.GetConfigDir(), logger.GetConfigName()))
			assert.Nil(t, err)
			// TEST
			assert.Equal(t, expectCode, w.Code)
			assert.Contains(t, string(logged), `"code":204`)
			assert.Contains(t, string(logged), expectMsg)
			os.Remove(path.Join(logger.GetConfigDir(), logger.GetConfigName()))
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

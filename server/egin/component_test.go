package egin

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/elog"
)

func TestComponent_buildTLSConfig(t *testing.T) {
	server := startClientAuthTLSServer()
	err := server.Init()
	assert.Nil(t, err)
	server.GET("/clientAuth", func(context *gin.Context) {
		context.String(200, "success")
	})
	go func() {
		assert.Nil(t, server.Start())
	}()
	defer func() {
		_ = server.Stop()
	}()
	time.Sleep(10 * time.Millisecond)
	t.Run("clientAuth", func(t *testing.T) {
		var c = &http.Client{
			Transport: &http.Transport{
				DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return tls.Dial(network, addr, loadConfig(t, true))
				},
			},
		}
		get, err1 := c.Get("https://127.0.0.1:20000/clientAuth")
		assert.Nil(t, err1)
		all, err1 := io.ReadAll(get.Body)
		assert.Nil(t, err1)
		assert.Equal(t, "success", string(all))
	})
	t.Run("NoClientAuth", func(t *testing.T) {
		var c = &http.Client{Transport: &http.Transport{
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return tls.Dial(network, addr, loadConfig(t, false))
			},
		}}
		_, err2 := c.Get("https://127.0.0.1:20000/clientAuth")
		assert.NotNil(t, err2)
	})

}

func TestContextClientIP(t *testing.T) {
	router := DefaultContainer().Build(WithTrustedPlatform("X-Forwarded-For"))
	router.GET("/", func(c *gin.Context) {
		assert.Equal(t, "10.10.10.11", c.ClientIP())
	})

	performRequest(router, "GET", "/", header{
		Key:   "X-Forwarded-For",
		Value: "10.10.10.11",
	})

	router3 := DefaultContainer().Build(WithTrustedPlatform("X-Forwarded-For"))
	router3.GET("/", func(c *gin.Context) {
		assert.NotEqual(t, "10.10.10.12", c.ClientIP())
	})

	performRequest(router3, "GET", "/", header{
		Key:   "X-Forwarded-For",
		Value: "10.10.10.11,10.10.10.12",
	})

	router2 := DefaultContainer().Build(WithTrustedPlatform("X-CUSTOM-CDN-IP"))
	router2.GET("/", func(c *gin.Context) {
		assert.Equal(t, "10.10.10.12", c.ClientIP())
	})

	performRequest(router2, "GET", "/", header{
		Key:   "X-CUSTOM-CDN-IP",
		Value: "10.10.10.12",
	})
}

func TestNewComponent(t *testing.T) {
	cfg := Config{
		Host:    "0.0.0.0",
		Port:    9006,
		Network: "tcp",
	}
	cmp := newComponent("test-cmp", &cfg, elog.DefaultLogger)
	assert.Equal(t, "test-cmp", cmp.Name())
	assert.Equal(t, "server.egin", cmp.PackageName())
	assert.Equal(t, "0.0.0.0:9006", cmp.config.Address())

	assert.NoError(t, cmp.Init())

	info := cmp.Info()
	assert.NotEmpty(t, info.Name)
	assert.Equal(t, "http", info.Scheme)
	assert.Equal(t, "[::]:9006", info.Address)
	assert.Equal(t, constant.ServiceProvider, info.Kind)

	// err = cmp.Start()
	go func() {
		assert.NoError(t, cmp.Start())
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	<-ctx.Done()
	assert.NoError(t, cmp.Stop())
}

func startClientAuthTLSServer() *Component {
	config := DefaultConfig()
	config.Port = 20000
	config.EnableTLS = true
	config.TLSKeyFile = "./testdata/egoServer/egoServer-key.pem"
	config.TLSCertFile = "./testdata/egoServer/egoServer.pem"
	config.TLSClientAuth = "RequireAndVerifyClientCert"
	config.TLSClientCAs = []string{"./testdata/egoClient/ca.pem"}
	container := DefaultContainer()
	container.config = config
	return container.Build()
}

func loadConfig(t *testing.T, loadClientCert bool) *tls.Config {
	pool := x509.NewCertPool()
	ca, err := os.ReadFile("./testdata/egoServer/ca.pem")
	assert.Nil(t, err)
	assert.True(t, pool.AppendCertsFromPEM(ca))
	cf := &tls.Config{}
	cf.RootCAs = pool
	if loadClientCert {
		serverCert, err := tls.LoadX509KeyPair("./testdata/egoClient/anyClient.pem", "./testdata/egoClient/anyClient-key.pem")
		assert.Nil(t, err)
		cf.Certificates = []tls.Certificate{serverCert}
	}
	return cf
}

func TestServerReadTimeout(t *testing.T) {
	// 使用非异步日志
	logger := elog.DefaultContainer().Build(
		elog.WithDebug(false),
		elog.WithEnableAddCaller(true),
		elog.WithEnableAsync(false),
	)
	timeout := 2 * time.Second
	container := DefaultContainer()
	cmp := container.Build(
		WithServerReadHeaderTimeout(2*time.Second),
		WithNetwork("local"),
		WithLogger(logger),
	)

	// 使用recover组件
	cmp.Use(container.defaultServerInterceptor())
	cmp.GET("/test", func(ctx *gin.Context) {
		time.Sleep(20 * time.Second)
		ctx.String(200, "hello world")
	})

	_ = cmp.Init()
	go func() {
		_ = cmp.Start()
	}()
	time.Sleep(1 * time.Second)

	// Slow client that should timeout.
	t1 := time.Now()
	conn, err := net.Dial("tcp", cmp.Listener().Addr().String())
	assert.Nil(t, err)

	buf := make([]byte, 1)
	n, err := conn.Read(buf)
	_ = conn.Close()
	latency := time.Since(t1)
	logger.Info("cost1", zap.Duration("cost", latency))
	if n != 0 || err != io.EOF {
		t.Error(fmt.Errorf("Read = %v, %v, wanted %v, %v", n, err, 0, io.EOF))
		return
	}
	minLatency := timeout / 5 * 4
	if latency < minLatency {
		t.Error(fmt.Errorf("got EOF after %s, want >= %s", latency, minLatency))
		return
	}
	fmt.Printf("path.Join(logger.ConfigDir(), logger.ConfigName())--------------->"+"%+v\n", path.Join(logger.ConfigDir(), logger.ConfigName()))
	_ = os.Remove(path.Join(logger.ConfigDir(), logger.ConfigName()))
}

func TestContextTimeout(t *testing.T) {
	// 使用非异步日志
	logger := elog.DefaultContainer().Build(
		elog.WithDebug(false),
		elog.WithEnableAddCaller(true),
		elog.WithEnableAsync(false),
	)
	timeout := 2 * time.Second
	container := DefaultContainer()
	cmp := container.Build(
		WithContextTimeout(timeout),
		WithNetwork("local"),
		WithLogger(logger),
	)

	// 使用recover组件
	cmp.Use(container.defaultServerInterceptor())
	cmp.GET("/test", func(ctx *gin.Context) {
		eginClient(ctx.Request.Context(), cmp, "/longTime")
		ctx.String(200, "hello world")
	})
	cmp.GET("/longTime", func(ctx *gin.Context) {
		time.Sleep(20 * time.Second)
		ctx.String(200, "i cost long time")
	})
	_ = cmp.Init()
	go func() {
		_ = cmp.Start()
	}()
	time.Sleep(1 * time.Second)

	// Slow client that should timeout.
	t1 := time.Now()
	err := eginClient(context.Background(), cmp, "/test")
	assert.Nil(t, err)

	latency := time.Since(t1)
	logger.Info("cost2", zap.Duration("cost", latency))
	os.Remove(path.Join(logger.ConfigDir(), logger.ConfigName()))
}

func TestServerTimeouts(t *testing.T) {
	timeout := 2 * time.Second
	err := testServerTimeouts(timeout)
	if err == nil {
		return
	}
	t.Logf("failed at %v: %v", timeout, err)
}

func testServerTimeouts(timeout time.Duration) error {
	reqNum := 0
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		reqNum++
		fmt.Fprintf(res, "req=%d", reqNum)
	}))
	ts.Config.ReadTimeout = timeout
	ts.Config.WriteTimeout = timeout
	ts.Start()
	defer ts.Close()

	// Hit the HTTP server successfully.
	c := ts.Client()
	t0 := time.Now()
	r, err := c.Get(ts.URL)
	if err != nil {
		return fmt.Errorf("http Get #1: %v", err)
	}
	got, err := io.ReadAll(r.Body)
	latency := time.Since(t0)
	fmt.Printf("got--------------->"+"%+v\n", got)
	fmt.Printf("latency--------------->"+"%+v\n", latency)

	expected := "req=1"
	if string(got) != expected || err != nil {
		return fmt.Errorf("Unexpected response for request #1; got %q ,%v; expected %q, nil",
			string(got), err, expected)
	}

	// Slow client that should timeout.
	t1 := time.Now()
	conn, err := net.Dial("tcp", ts.Listener.Addr().String())
	if err != nil {
		return fmt.Errorf("Dial: %v", err)
	}
	buf := make([]byte, 1)
	n, err := conn.Read(buf)
	conn.Close()
	latency = time.Since(t1)
	fmt.Printf("latency--------------->"+"%+v\n", latency)
	if n != 0 || err != io.EOF {
		return fmt.Errorf("Read = %v, %v, wanted %v, %v", n, err, 0, io.EOF)
	}
	minLatency := timeout / 5 * 4
	if latency < minLatency {
		return fmt.Errorf("got EOF after %s, want >= %s", latency, minLatency)
	}

	// Hit the HTTP server successfully again, verifying that the
	// previous slow connection didn't run our handler.  (that we
	// get "req=2", not "req=3")
	r, err = c.Get(ts.URL)
	if err != nil {
		return fmt.Errorf("http Get #2: %v", err)
	}
	got, err = io.ReadAll(r.Body)
	r.Body.Close()
	expected = "req=2"

	if string(got) != expected || err != nil {
		return fmt.Errorf("Get #2 got %q, %v, want %q, nil", string(got), err, expected)
	}
	return nil
}

func eginClient(ctx context.Context, gin *Component, url string) (err error) {
	client := &http.Client{Transport: &http.Transport{}}
	req, err := http.NewRequestWithContext(ctx, "GET", "http://"+gin.Listener().Addr().String()+url, nil)
	if err != nil {
		return err
	}
	r, err := client.Do(req)
	_, err = io.ReadAll(r.Body)
	r.Body.Close()
	return nil
}

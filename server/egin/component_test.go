package egin

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/elog"
	"github.com/stretchr/testify/assert"
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
		Host: "0.0.0.0",
		Port: 9005,
	}
	cmp := newComponent("test-cmp", &cfg, elog.DefaultLogger)
	assert.Equal(t, "test-cmp", cmp.Name())
	assert.Equal(t, "server.egin", cmp.PackageName())
	assert.Equal(t, "0.0.0.0:9005", cmp.config.Address())

	assert.NoError(t, cmp.Init())

	info := cmp.Info()
	assert.NotEmpty(t, info.Name)
	assert.Equal(t, "http", info.Scheme)
	assert.Equal(t, "[::]:9005", info.Address)
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

package egin

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestComponent_buildTLSConfig(t *testing.T) {
	server := startClientAuthTLSServer()
	err := server.Init()
	assert.Nil(t, err)
	server.GET("/clientAuth", func(context *gin.Context) {
		context.String(200, "success")
	})
	go func() {
		err := server.Start()
		assert.Nil(t, err)
	}()
	defer server.Stop()
	time.Sleep(10 * time.Millisecond)
	t.Run("clientAuth", func(t *testing.T) {
		var c = &http.Client{
			Transport: &http.Transport{
				DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return tls.Dial(network, addr, loadConfig(t, true))
				},
			},
		}
		get, err := c.Get("https://127.0.0.1:20000/clientAuth")
		assert.Nil(t, err)
		all, err := io.ReadAll(get.Body)
		assert.Nil(t, err)
		assert.Equal(t, "success", string(all))
	})
	t.Run("NoClientAuth", func(t *testing.T) {
		var c = &http.Client{Transport: &http.Transport{
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return tls.Dial(network, addr, loadConfig(t, false))
			},
		}}
		_, err := c.Get("https://127.0.0.1:20000/clientAuth")
		assert.NotNil(t, err)
	})

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

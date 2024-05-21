package egin

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/gotomicro/ego/core/util/xtime"
)

func TestDefaultConfig(t *testing.T) {
	in := &Config{
		// Host:                       eflag.String("host"),
		Port:                       9090,
		Mode:                       gin.ReleaseMode,
		Network:                    "tcp",
		EnableAccessInterceptor:    true,
		EnableTraceInterceptor:     true,
		EnableMetricInterceptor:    true,
		EnableSentinel:             true,
		SlowLogThreshold:           xtime.Duration("500ms"),
		EnableWebsocketCheckOrigin: false,
		TrustedPlatform:            "",
		recoveryFunc:               defaultRecoveryFunc,
	}
	out := DefaultConfig()
	assert.Equal(t, in.Port, out.Port)
	assert.Equal(t, in.Mode, out.Mode)
	assert.Equal(t, in.Network, out.Network)
	assert.Equal(t, in.EnableAccessInterceptor, out.EnableAccessInterceptor)
	assert.Equal(t, in.EnableTraceInterceptor, out.EnableTraceInterceptor)
	assert.Equal(t, in.EnableMetricInterceptor, out.EnableMetricInterceptor)
	assert.Equal(t, in.EnableSentinel, out.EnableSentinel)
	assert.Equal(t, in.SlowLogThreshold, out.SlowLogThreshold)
	assert.Equal(t, in.EnableWebsocketCheckOrigin, out.EnableWebsocketCheckOrigin)
	assert.Equal(t, in.TrustedPlatform, out.TrustedPlatform)
}

func TestAddress(t *testing.T) {
	config := Config{
		Host: PackageName,
		Port: 9090,
	}
	out := config.Address()
	assert.Equal(t, "server.egin:9090", out)
}

func TestClientAuthType(t *testing.T) {
	config := &Config{TLSClientAuth: "RequireAnyClientCert"}
	assert.Equal(t, tls.RequireAnyClientCert, config.ClientAuthType())

	config.TLSClientAuth = "RequestClientCert"
	assert.Equal(t, tls.RequestClientCert, config.ClientAuthType())

	config.TLSClientAuth = "VerifyClientCertIfGiven"
	assert.Equal(t, tls.VerifyClientCertIfGiven, config.ClientAuthType())

	config.TLSClientAuth = "RequireAndVerifyClientCert"
	assert.Equal(t, tls.RequireAndVerifyClientCert, config.ClientAuthType())

	config.TLSClientAuth = "NoClientCert"
	assert.Equal(t, tls.NoClientCert, config.ClientAuthType())
}

func TestDefaultRecoveryFunc(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusInternalServerError)
	})
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	fmt.Printf("w.Code: %v\n", w.Code)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

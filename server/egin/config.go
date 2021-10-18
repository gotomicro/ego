package egin

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gotomicro/ego/core/eflag"
	"github.com/gotomicro/ego/core/util/xtime"
)

// Config HTTP config
type Config struct {
	Host                       string        // IP地址，默认0.0.0.0
	Port                       int           // PORT端口，默认9001
	Mode                       string        // gin的模式，默认是release模式
	EnableMetricInterceptor    bool          // 是否开启监控，默认开启
	EnableTraceInterceptor     bool          // 是否开启链路追踪，默认开启
	EnableLocalMainIP          bool          // 自动获取ip地址
	SlowLogThreshold           time.Duration // 服务慢日志，默认500ms
	EnableAccessInterceptor    bool          // 是否开启，记录请求数据
	EnableAccessInterceptorReq bool          // 是否开启记录请求参数，默认不开启
	EnableAccessInterceptorRes bool          // 是否开启记录响应参数，默认不开启
	EnableTrustedCustomHeader  bool          // 是否开启自定义header头，记录数据往链路后传递，默认不开启
	EnableSentinel             bool          // 是否开启限流，默认不开启
	WebsocketHandshakeTimeout  time.Duration // 握手时间
	WebsocketReadBufferSize    int
	WebsocketWriteBufferSize   int
	EnableWebsocketCompression bool     // 是否开通压缩
	EnableWebsocketCheckOrigin bool     // 是否支持跨域
	EnableTLS                  bool     // 是否进入 https 模式
	TLSCertFile                string   // https 证书
	TLSKeyFile                 string   // https 私钥
	TLSClientAuth              string   // https 客户端认证方式默认为 NoClientCert(NoClientCert,RequestClientCert,RequireAnyClientCert,VerifyClientCertIfGiven,RequireAndVerifyClientCert)
	TLSClientCAs               []string // https client的ca，当需要双向认证的时候指定可以倒入自签证书
	TLSSessionCache            tls.ClientSessionCache
	blockFallback              func(*gin.Context)
	resourceExtract            func(*gin.Context) string
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Host:                       eflag.String("host"),
		Port:                       9090,
		Mode:                       gin.ReleaseMode,
		EnableAccessInterceptor:    true,
		EnableTraceInterceptor:     true,
		EnableMetricInterceptor:    true,
		SlowLogThreshold:           xtime.Duration("500ms"),
		EnableWebsocketCheckOrigin: false,
	}
}

// Address ...
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

// ClientAuthType 客户端auth类型
func (config *Config) ClientAuthType() tls.ClientAuthType {
	switch config.TLSClientAuth {
	case "NoClientCert":
		return tls.NoClientCert
	case "RequestClientCert":
		return tls.RequestClientCert
	case "RequireAnyClientCert":
		return tls.RequireAnyClientCert
	case "VerifyClientCertIfGiven":
		return tls.VerifyClientCertIfGiven
	case "RequireAndVerifyClientCert":
		return tls.RequireAndVerifyClientCert
	default:
		return tls.NoClientCert
	}
}

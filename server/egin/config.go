package egin

import (
	"crypto/tls"
	"embed"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/cel-go/cel"

	"github.com/gotomicro/ego/core/eflag"
	"github.com/gotomicro/ego/core/util/xtime"
)

// Config HTTP config
type Config struct {
	Host                    string // IP地址，默认0.0.0.0
	Port                    int    // PORT端口，默认9001
	Mode                    string // gin的模式，默认是release模式
	Network                 string
	ServerReadTimeout       time.Duration // 服务端，用于读取io报文过慢的timeout，通常用于互联网网络收包过慢，如果你的go在最外层，可以使用他，默认不启用。
	ServerReadHeaderTimeout time.Duration // 服务端，用于读取io报文过慢的timeout，通常用于互联网网络收包过慢，如果你的go在最外层，可以使用他，默认不启用。
	ServerWriteTimeout      time.Duration // 服务端，用于读取io报文过慢的timeout，通常用于互联网网络收包过慢，如果你的go在最外层，可以使用他，默认不启用。
	// ServerHTTPTimout        time.Duration //  这个是HTTP包提供的，可以用于IO，或者密集型计算，做timeout处理，有一次goroutine操作，然后没走一些流程，cancel体验不好，暂时先不用
	ContextTimeout                time.Duration // 只能用于IO操作，才能触发，默认不启用
	EnableMetricInterceptor       bool          // 是否开启监控，默认开启
	EnableTraceInterceptor        bool          // 是否开启链路追踪，默认开启
	EnableLocalMainIP             bool          // 自动获取ip地址
	SlowLogThreshold              time.Duration // 服务慢日志，默认500ms
	EnableAccessInterceptor       bool          // 是否开启，记录请求数据
	EnableAccessInterceptorReq    bool          // 是否开启记录请求参数，默认不开启
	EnableAccessInterceptorRes    bool          // 是否开启记录响应参数，默认不开启
	AccessInterceptorReqResFilter string        // AccessInterceptorReq 过滤器，只有符合过滤器的请求才会记录 Req 和 Res
	EnableTrustedCustomHeader     bool          // 是否开启自定义header头，记录数据往链路后传递，默认不开启
	EnableSentinel                bool          // 是否开启限流，默认不开启
	WebsocketHandshakeTimeout     time.Duration // 握手时间
	WebsocketReadBufferSize       int           // WebsocketReadBufferSize
	WebsocketWriteBufferSize      int           // WebsocketWriteBufferSize
	EnableWebsocketCompression    bool          // 是否开通压缩
	EnableWebsocketCheckOrigin    bool          // 是否支持跨域
	EnableTLS                     bool          // 是否进入 https 模式
	TLSCertFile                   string        // https 证书
	TLSKeyFile                    string        // https 私钥
	TLSClientAuth                 string        // https 客户端认证方式默认为 NoClientCert(NoClientCert,RequestClientCert,RequireAnyClientCert,VerifyClientCertIfGiven,RequireAndVerifyClientCert)
	TLSClientCAs                  []string      // https client的ca，当需要双向认证的时候指定可以倒入自签证书
	TrustedPlatform               string        // 需要用户换成自己的CDN名字，获取客户端IP地址
	EmbedPath                     string        // 嵌入embed path数据
	EnableH2C                     bool          // 开启HTTP2
	embedFs                       embed.FS      // 需要在build时候注入embed.Fs
	TLSSessionCache               tls.ClientSessionCache
	blockFallback                 func(*gin.Context)
	resourceExtract               func(*gin.Context) string
	aiReqResCelPrg                cel.Program
	mu                            sync.RWMutex     // mutex for EnableAccessInterceptorReq、EnableAccessInterceptorRes、AccessInterceptorReqResFilter、aiReqResCelPrg
	recoveryFunc                  gin.RecoveryFunc // recoveryFunc 处理接口没有被 recover 的 panic，默认返回 500 并且没有任何 response body
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Host:                       eflag.String("host"),
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

func defaultRecoveryFunc(ctx *gin.Context, _ interface{}) {
	ctx.AbortWithStatus(http.StatusInternalServerError)
}

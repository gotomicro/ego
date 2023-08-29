package egin

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	healthcheck "github.com/RaMin0/gin-health-check"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/eapp"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server"
)

// PackageName 包名
const PackageName = "server.egin"

// Component ...
type Component struct {
	*gin.Engine
	mu               sync.Mutex
	name             string
	config           *Config
	logger           *elog.Component
	Server           *http.Server
	listener         net.Listener
	routerCommentMap map[string]string // router的中文注释，非并发安全
	embedWrapper     *EmbedWrapper
	invokers         []func() error // 用户初始化函数
}

func newComponent(name string, config *Config, logger *elog.Component) *Component {
	gin.SetMode(config.Mode)
	comp := &Component{
		name:             name,
		config:           config,
		logger:           logger,
		Engine:           gin.New(),
		listener:         nil,
		routerCommentMap: make(map[string]string),
	}

	if config.EmbedPath != "" {
		comp.embedWrapper = &EmbedWrapper{
			embedFs: config.embedFs,
			path:    config.EmbedPath,
		}
	}

	// 设置信任的header头
	comp.Engine.TrustedPlatform = config.TrustedPlatform
	return comp
}

// Name 配置名称
func (c *Component) Name() string {
	return c.name
}

// PackageName 包名
func (c *Component) PackageName() string {
	return PackageName
}

// Prepare 预准备数据
func (c *Component) Prepare() error {
	for _, fn := range c.invokers {
		err := fn()
		if err != nil {
			return err
		}
	}
	return nil
}

// Init 初始化
func (c *Component) Init() error {
	var err error
	if c.config.Network == "local" {
		c.listener = newLocalListener()
	} else {
		c.listener, err = net.Listen(c.config.Network, c.config.Address())
		if err != nil {
			c.logger.Panic("new egin server err", elog.FieldErrKind("listen err"), elog.FieldErr(err))
		}
	}
	c.config.Port = c.listener.Addr().(*net.TCPAddr).Port
	return nil
}

// Health implements server.Component interface.
func (c *Component) Health() bool {
	cli, err := resty.New().
		SetHeader("app", eapp.Name()).
		SetBaseURL("http://"+c.config.Address()).R().SetHeader(healthcheck.DefaultHeaderName, healthcheck.DefaultHeaderValue).Get("/")
	if err != nil {
		c.logger.Error("health fail", elog.FieldErr(err))
		return false
	}
	return cli.StatusCode() == healthcheck.DefaultResponseCode && cli.String() == healthcheck.DefaultResponseText
}

// RegisterRouteComment 注册路由注释
func (c *Component) RegisterRouteComment(method, path, comment string) {
	c.routerCommentMap[commentUniqKey(method, path)] = comment
}

// Start implements server.Component interface.
func (c *Component) Start() error {
	var err error
	for _, route := range c.Engine.Routes() {
		info, flag := c.routerCommentMap[commentUniqKey(route.Method, route.Path)]
		// 如果有注释，日志打出来
		if flag {
			c.logger.Info("add route", elog.FieldMethod(route.Method), elog.String("path", route.Path), elog.Any("info", info))
		} else {
			c.logger.Info("add route", elog.FieldMethod(route.Method), elog.String("path", route.Path))
		}
	}

	// 因为start和stop在多个goroutine里，需要对Server上写锁
	c.mu.Lock()
	c.Server = &http.Server{
		Addr:              c.config.Address(),
		Handler:           c,
		ReadHeaderTimeout: c.config.ServerReadHeaderTimeout,
		ReadTimeout:       c.config.ServerReadTimeout,
		WriteTimeout:      c.config.ServerWriteTimeout,
		// Handler:           http.TimeoutHandler(c, 1*time.Second, "timeout"),
	}
	c.mu.Unlock()
	if c.config.EnableTLS {
		config, errTLS := c.buildTLSConfig()
		if errTLS != nil {
			return errTLS
		}
		c.Server.TLSConfig = config
		err = c.Server.ServeTLS(c.listener, "", "")
	} else {
		err = c.Server.Serve(c.listener)
	}
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

// Stop implements server.Component interface
// it will terminate gin server immediately
func (c *Component) Stop() error {
	c.mu.Lock()
	err := c.Server.Close()
	c.mu.Unlock()
	return err
}

// GracefulStop implements server.Component interface
// it will stop gin server gracefully
func (c *Component) GracefulStop(ctx context.Context) error {
	c.mu.Lock()
	err := c.Server.Shutdown(ctx)
	c.mu.Unlock()
	return err
}

// Info returns server info, used by governor and consumer balancer
func (c *Component) Info() *server.ServiceInfo {
	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(c.listener.Addr().String()),
		server.WithKind(constant.ServiceProvider),
	)
	return &info
}

// Invoker returns server info, used by governor and consumer balancer
// Experimental
func (c *Component) Invoker(fns ...func() error) {
	c.invokers = append(c.invokers, fns...)
}

func commentUniqKey(method, path string) string {
	return fmt.Sprintf("%s@%s", strings.ToLower(method), path)
}

func (c *Component) buildTLSConfig() (*tls.Config, error) {
	tlsConfig := &tls.Config{}
	serverCert, err := tls.LoadX509KeyPair(c.config.TLSCertFile, c.config.TLSKeyFile)
	if err != nil {
		return nil, err
	}
	tlsConfig.Certificates = []tls.Certificate{serverCert}
	tlsConfig.ClientCAs = x509.NewCertPool()
	tlsConfig.ClientAuth = c.config.ClientAuthType()
	tlsConfig.ClientSessionCache = c.config.TLSSessionCache
	clientCAs := c.config.TLSClientCAs
	for i := range clientCAs {
		clientCA := clientCAs[i]
		ca, err := os.ReadFile(clientCA)
		if err != nil {
			return nil, fmt.Errorf("read client ca fail:%+v", err)
		}
		if !tlsConfig.ClientCAs.AppendCertsFromPEM(ca) {
			return nil, fmt.Errorf("append client ca fail:%+v", err)
		}
	}
	return tlsConfig, nil
}

// HTTPEmbedFs http的文件系统
func (c *Component) HTTPEmbedFs() http.FileSystem {
	return http.FS(c.embedWrapper)
}

// GetEmbedWrapper http的文件系统
func (c *Component) GetEmbedWrapper() *EmbedWrapper {
	return c.embedWrapper
}

// Listener listener信息
func (c *Component) Listener() net.Listener {
	return c.listener
}

// EmbedWrapper 嵌入普通的静态资源的wrapper
type EmbedWrapper struct {
	embedFs embed.FS // 静态资源
	path    string   // 设置embed文件到静态资源的相对路径，也就是embed注释里的路径
}

// Open 静态资源被访问的核心逻辑
func (e *EmbedWrapper) Open(name string) (fs.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	fullName := filepath.ToSlash(path.Join(e.path, path.Clean("/"+name)))
	file, err := e.embedFs.Open(fullName)
	return file, err
}

func newLocalListener() net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if l, err = net.Listen("tcp6", "[::1]:0"); err != nil {
			panic(fmt.Sprintf("httptest: failed to listen on a port: %v", err))
		}
	}
	return l
}


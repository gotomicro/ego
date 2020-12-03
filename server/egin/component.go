package egin

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server"
)

const PackageName = "server.egin"

// Component ...
type Component struct {
	name   string
	config *Config
	logger *elog.Component
	*gin.Engine
	Server           *http.Server
	listener         net.Listener
	routerCommentMap map[string]string // router的中文注释，非并发安全
}

func newComponent(name string, config *Config, logger *elog.Component) *Component {
	gin.SetMode(config.Mode)
	return &Component{
		name:             name,
		config:           config,
		logger:           logger,
		Engine:           gin.New(),
		listener:         nil,
		routerCommentMap: make(map[string]string),
	}
}

func (c *Component) Name() string {
	return c.name
}

func (c *Component) PackageName() string {
	return PackageName
}

func (c *Component) Init() error {
	listener, err := net.Listen("tcp", c.config.Address())
	if err != nil {
		c.logger.Panic("new egin server err", elog.FieldErrKind("listen err"), elog.FieldErr(err))
	}
	c.config.Port = listener.Addr().(*net.TCPAddr).Port
	c.listener = listener
	return nil
}

// 注册路由注释
func (c *Component) RegisterRouteComment(method, path, comment string) {
	c.routerCommentMap[commentUniqKey(method, path)] = comment
}

//Upgrade protocol to WebSocket
func (c *Component) Upgrade(ws *WebSocket) gin.IRoutes {
	return c.GET(ws.Pattern, func(c *gin.Context) {
		ws.Upgrade(c.Writer, c.Request)
	})
}

// Serve implements server.Component interface.
func (c *Component) Start() error {

	for _, route := range c.Engine.Routes() {
		info, flag := c.routerCommentMap[commentUniqKey(route.Method, route.Path)]
		// 如果有注释，日志打出来
		if flag {
			c.logger.Info("add route", elog.FieldMethod(route.Method), elog.String("path", route.Path), elog.Any("info", info))
		} else {
			c.logger.Info("add route", elog.FieldMethod(route.Method), elog.String("path", route.Path))
		}
	}
	c.Server = &http.Server{
		Addr:    c.config.Address(),
		Handler: c,
	}
	err := c.Server.Serve(c.listener)
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

// Stop implements server.Component interface
// it will terminate gin server immediately
func (c *Component) Stop() error {
	return c.Server.Close()
}

// GracefulStop implements server.Component interface
// it will stop gin server gracefully
func (c *Component) GracefulStop(ctx context.Context) error {
	return c.Server.Shutdown(ctx)
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

func commentUniqKey(method, path string) string {
	return fmt.Sprintf("%s@%s", strings.ToLower(method), path)
}

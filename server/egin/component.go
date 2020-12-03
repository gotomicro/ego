package egin

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

//WebSocketConn websocket conn, see websocket.Conn
type WebSocketConn interface {
	Subprotocol() string
	Close() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	WriteControl(messageType int, data []byte, deadline time.Time) error
	NextWriter(messageType int) (io.WriteCloser, error)
	WritePreparedMessage(pm *websocket.PreparedMessage) error
	WriteMessage(messageType int, data []byte) error
	SetWriteDeadline(t time.Time) error
	NextReader() (messageType int, r io.Reader, err error)
	ReadMessage() (messageType int, p []byte, err error)
	SetReadDeadline(t time.Time) error
	SetReadLimit(limit int64)
	CloseHandler() func(code int, text string) error
	SetCloseHandler(h func(code int, text string) error)
	PingHandler() func(appData string) error
	SetPingHandler(h func(appData string) error)
	PongHandler() func(appData string) error
	SetPongHandler(h func(appData string) error)
	UnderlyingConn() net.Conn
	EnableWriteCompression(enable bool)
	SetCompressionLevel(level int) error
}

//WebSocketFunc ..
type WebSocketFunc func(WebSocketConn, error)

//WebSocket ..
type WebSocket struct {
	Pattern string
	Handler WebSocketFunc
	*websocket.Upgrader
	Header http.Header
}

//Upgrade get upgrage request
func (ws *WebSocket) Upgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.Upgrader.Upgrade(w, r, ws.Header)
	if err == nil {
		defer conn.Close()
	}
	ws.Handler(conn, err)
}

//WebSocketOption ..
type WebSocketOption func(*WebSocket)

//WebSocketOptions ..
func WebSocketOptions(pattern string, handler WebSocketFunc, opts ...WebSocketOption) *WebSocket {
	ws := &WebSocket{
		Pattern:  pattern,
		Handler:  handler,
		Upgrader: &websocket.Upgrader{},
	}
	for _, opt := range opts {
		opt(ws)
	}
	return ws
}

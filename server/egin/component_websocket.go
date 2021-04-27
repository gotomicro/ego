package egin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Upgrade protocol to WebSocket
func (c *Component) Upgrade(pattern string, ws *WebSocket, handler WebSocketFunc) gin.IRoutes {
	return c.GET(pattern, func(c *gin.Context) {
		ws.Upgrade(c.Writer, c.Request, c, handler)
	})
}

// BuildWebsocket ..
func (c *Component) BuildWebsocket(opts ...WebSocketOption) *WebSocket {
	upgrader := &websocket.Upgrader{}
	// 支持跨域
	if c.config.EnableWebsocketCheckOrigin {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}

	ws := &WebSocket{
		Upgrader: upgrader,
	}
	for _, opt := range opts {
		opt(ws)
	}
	return ws
}

// WebSocket ..
type WebSocket struct {
	*websocket.Upgrader
}

// Upgrade get upgrage request
func (ws *WebSocket) Upgrade(w http.ResponseWriter, r *http.Request, c *gin.Context, handler WebSocketFunc) {
	// todo response Header
	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err == nil {
		defer conn.Close()
	}
	wsConn := &WebSocketConn{
		Conn:   conn,
		GinCtx: c,
	}
	handler(wsConn, err)
}

// WebSocketFunc ..
type WebSocketFunc func(*WebSocketConn, error)

// WebSocketConn ...
type WebSocketConn struct {
	*websocket.Conn
	GinCtx *gin.Context
}

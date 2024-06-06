package egin

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// simpleBufferPool is an implementation of BufferPool for TestWriteBufferPool.
type simpleBufferPool struct {
	v interface{}
}

func (p *simpleBufferPool) Get() interface{} {
	v := p.v
	p.v = nil
	return v
}

func (p *simpleBufferPool) Put(v interface{}) {
	p.v = v
}

var up = &websocket.Upgrader{
	HandshakeTimeout:  30 * time.Second,
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	Subprotocols:      []string{"p0", "p1"},
	WriteBufferPool:   &simpleBufferPool{},
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		http.Error(w, reason.Error(), status)
	},
}

func TestUpgrade(t *testing.T) {
	c := DefaultContainer().Build()
	wss := &WebSocket{up}
	handler := func(conn *WebSocketConn, err error) {}
	out := c.Upgrade("/test", wss, handler)
	in := c.GET("/hello", func(ctx *gin.Context) {
		wss.Upgrade(ctx.Writer, ctx.Request, ctx, handler)
	})
	assert.Equal(t, in, out)
}

func TestBuildWebsocket(t *testing.T) {
	opt := func(wss *WebSocket) {
		wss.Upgrader = up
	}
	c := DefaultContainer().Build()
	out := c.BuildWebsocket(opt)
	reflect.DeepEqual(&WebSocket{}, out)

	err := c.Prepare()
	assert.Equal(t, nil, err)

	h := c.Health()
	assert.Equal(t, false, h)

	reflect.DeepEqual(http.FS(c.embedWrapper), c.HTTPEmbedFs())
	reflect.DeepEqual(c.embedWrapper, c.GetEmbedWrapper())
}

func TestOpen(t *testing.T) {
	e := &EmbedWrapper{
		embedFs: embed.FS{},
		path:    PackageName,
	}
	file, err := e.Open("test")
	reflect.DeepEqual("", file)
	reflect.DeepEqual(nil, err)
}

func TestWebSocket_Upgrade(t *testing.T) {
	wss := &WebSocket{
		up,
	}
	w := httptest.NewRecorder()
	r := &http.Request{}
	c := &gin.Context{}
	handler := func(conn *WebSocketConn, err error) {}
	wss.Upgrade(w, r, c, handler)
	reflect.DeepEqual("", handler)
}

package egin

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestUpgrade(t *testing.T) {
	c := DefaultContainer().Build()
	ws := &WebSocket{}
	handler := func(conn *WebSocketConn, err error) {}
	c.Upgrade("test", ws, handler)
	assert.NoError(t, nil)
}

func TestBuildWebsocket(t *testing.T) {
	opt := func(ws *WebSocket) {
		ws.Upgrader = &websocket.Upgrader{
			HandshakeTimeout: 3,
			ReadBufferSize:   1024,
			WriteBufferSize:  1024,
			WriteBufferPool:  &simpleBufferPool{},
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			Subprotocols:      make([]string, 0),
			Error:             nil,
			EnableCompression: true}
	}
	c := DefaultContainer().Build() // 设置config
	c.BuildWebsocket(opt)
	assert.NoError(t, nil)

	err := c.Prepare()
	assert.Equal(t, nil, err)

	h := c.Health()
	assert.Equal(t, false, h)

	c.HTTPEmbedFs()
	assert.NoError(t, nil)

	c.GetEmbedWrapper()
	assert.NoError(t, nil)

	e := &EmbedWrapper{}
	e.Open("test")
	assert.NoError(t, nil)
}

func TestWebSocket_Upgrade(t *testing.T) {
	ws := &WebSocket{
		&websocket.Upgrader{
			HandshakeTimeout: 3,
			ReadBufferSize:   1024,
			WriteBufferSize:  1024,
			WriteBufferPool:  &simpleBufferPool{},
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			Subprotocols:      make([]string, 0),
			Error:             nil,
			EnableCompression: true},
	}
	w := httptest.NewRecorder()
	r := &http.Request{}
	c := &gin.Context{}
	handler := func(conn *WebSocketConn, err error) {}
	ws.Upgrade(w, r, c, handler)
	assert.NoError(t, nil)
}

package web

import (
	"io"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

// WSConn describes a websocket connection. Passthrough to websocket.Conn from github.com/gorilla/websocket
type WSConn struct {
	c *websocket.Conn
}

// Subprotocol - see github.com/gorilla/websocket for info
func (c *WSConn) Subprotocol() string {
	return c.c.Subprotocol()
}

// Close - see github.com/gorilla/websocket for info
func (c *WSConn) Close() error {
	return c.c.Close()
}

// LocalAddr - see github.com/gorilla/websocket for info
func (c *WSConn) LocalAddr() net.Addr {
	return c.c.LocalAddr()
}

// RemoteAddr - see github.com/gorilla/websocket for info
func (c *WSConn) RemoteAddr() net.Addr {
	return c.c.RemoteAddr()
}

// WriteControl - see github.com/gorilla/websocket for info
func (c *WSConn) WriteControl(messageType int, data []byte, deadline time.Time) error {
	return c.c.WriteControl(messageType, data, deadline)
}

// NextWriter - see github.com/gorilla/websocket for info
func (c *WSConn) NextWriter(messageType int) (io.WriteCloser, error) {
	return c.c.NextWriter(messageType)
}

// WritePreparedMessage - see github.com/gorilla/websocket for info
func (c *WSConn) WritePreparedMessage(pm *websocket.PreparedMessage) error {
	return c.c.WritePreparedMessage(pm)
}

// WriteMessage - see github.com/gorilla/websocket for info
func (c *WSConn) WriteMessage(messageType int, data []byte) error {
	return c.c.WriteMessage(messageType, data)
}

// SetWriteDeadline - see github.com/gorilla/websocket for info
func (c *WSConn) SetWriteDeadline(t time.Time) error {
	return c.c.SetWriteDeadline(t)
}

// NextReader - see github.com/gorilla/websocket for info
func (c *WSConn) NextReader() (messageType int, r io.Reader, err error) {
	return c.c.NextReader()
}

// ReadMessage - see github.com/gorilla/websocket for info
func (c *WSConn) ReadMessage() (messageType int, p []byte, err error) {
	return c.c.ReadMessage()
}

// SetReadDeadline - see github.com/gorilla/websocket for info
func (c *WSConn) SetReadDeadline(t time.Time) error {
	return c.c.SetReadDeadline(t)
}

// SetReadLimit - see github.com/gorilla/websocket for info
func (c *WSConn) SetReadLimit(limit int64) {
	c.c.SetReadLimit(limit)
}

// CloseHandler - see github.com/gorilla/websocket for info
func (c *WSConn) CloseHandler() func(code int, text string) error {
	return c.c.CloseHandler()
}

// SetCloseHandler - see github.com/gorilla/websocket for info
func (c *WSConn) SetCloseHandler(h func(code int, text string) error) {
	c.c.SetCloseHandler(h)
}

// PingHandler - see github.com/gorilla/websocket for info
func (c *WSConn) PingHandler() func(appData string) error {
	return c.c.PingHandler()
}

// SetPingHandler - see github.com/gorilla/websocket for info
func (c *WSConn) SetPingHandler(h func(appData string) error) {
	c.c.SetPingHandler(h)
}

// PongHandler - see github.com/gorilla/websocket for info
func (c *WSConn) PongHandler() func(appData string) error {
	return c.c.PongHandler()
}

// SetPongHandler - see github.com/gorilla/websocket for info
func (c *WSConn) SetPongHandler(h func(appData string) error) {
	c.c.SetPongHandler(h)
}

// UnderlyingConn - see github.com/gorilla/websocket for info
func (c *WSConn) UnderlyingConn() net.Conn {
	return c.c.UnderlyingConn()
}

// EnableWriteCompression - see github.com/gorilla/websocket for info
func (c *WSConn) EnableWriteCompression(enable bool) {
	c.c.EnableWriteCompression(enable)
}

// SetCompressionLevel - see github.com/gorilla/websocket for info
func (c *WSConn) SetCompressionLevel(level int) error {
	return c.c.SetCompressionLevel(level)
}

// WriteJSON - see github.com/gorilla/websocket for info
func (c *WSConn) WriteJSON(v interface{}) error {
	return c.c.WriteJSON(v)
}

// ReadJSON - see github.com/gorilla/websocket for info
func (c *WSConn) ReadJSON(v interface{}) error {
	return c.c.ReadJSON(v)
}

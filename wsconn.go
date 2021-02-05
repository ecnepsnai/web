package web

import (
	"io"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

// WSConn describes a websocket connection.
type WSConn struct {
	c *websocket.Conn
}

// Subprotocol returns the negotiated protocol for the connection.
func (c *WSConn) Subprotocol() string {
	return c.c.Subprotocol()
}

// Close closes the underlying network connection without sending or waiting for a close message.
func (c *WSConn) Close() error {
	return c.c.Close()
}

// LocalAddr returns the local network address.
func (c *WSConn) LocalAddr() net.Addr {
	return c.c.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *WSConn) RemoteAddr() net.Addr {
	return c.c.RemoteAddr()
}

// WriteControl writes a control message with the given deadline. The allowed message types are CloseMessage,
// PingMessage and PongMessage.
func (c *WSConn) WriteControl(messageType int, data []byte, deadline time.Time) error {
	return c.c.WriteControl(messageType, data, deadline)
}

// NextWriter returns a writer for the next message to send. The writer's Close
// method flushes the complete message to the network.
//
// There can be at most one open writer on a connection. NextWriter closes the
// previous writer if the application has not already done so.
//
// All message types (TextMessage, BinaryMessage, CloseMessage, PingMessage and
// PongMessage) are supported.
func (c *WSConn) NextWriter(messageType int) (io.WriteCloser, error) {
	return c.c.NextWriter(messageType)
}

// WritePreparedMessage writes prepared message into connection.
func (c *WSConn) WritePreparedMessage(pm *websocket.PreparedMessage) error {
	return c.c.WritePreparedMessage(pm)
}

// WriteMessage is a helper method for getting a writer using NextWriter, writing the message and closing the writer.
func (c *WSConn) WriteMessage(messageType int, data []byte) error {
	return c.c.WriteMessage(messageType, data)
}

// SetWriteDeadline sets the write deadline on the underlying network
// connection. After a write has timed out, the websocket state is corrupt and
// all future writes will return an error. A zero value for t means writes will
// not time out.
func (c *WSConn) SetWriteDeadline(t time.Time) error {
	return c.c.SetWriteDeadline(t)
}

// NextReader returns the next data message received from the peer. The
// returned messageType is either TextMessage or BinaryMessage.
//
// There can be at most one open reader on a connection. NextReader discards
// the previous message if the application has not already consumed it.
//
// Applications must break out of the application's read loop when this method
// returns a non-nil error value. Errors returned from this method are
// permanent. Once this method returns a non-nil error, all subsequent calls to
// this method return the same error.
func (c *WSConn) NextReader() (messageType int, r io.Reader, err error) {
	return c.c.NextReader()
}

// ReadMessage is a helper method for getting a reader using NextReader and reading from that reader to a buffer.
func (c *WSConn) ReadMessage() (messageType int, p []byte, err error) {
	return c.c.ReadMessage()
}

// SetReadDeadline sets the read deadline on the underlying network connection.
// After a read has timed out, the websocket connection state is corrupt and all future reads will return an error.
// A zero value for t means reads will not time out.
func (c *WSConn) SetReadDeadline(t time.Time) error {
	return c.c.SetReadDeadline(t)
}

// SetReadLimit sets the maximum size in bytes for a message read from the peer. If a message exceeds the limit, the
// connection sends a close message to the peer and returns ErrReadLimit to the application.
func (c *WSConn) SetReadLimit(limit int64) {
	c.c.SetReadLimit(limit)
}

// CloseHandler returns the current close handler
func (c *WSConn) CloseHandler() func(code int, text string) error {
	return c.c.CloseHandler()
}

// SetCloseHandler sets the handler for close messages received from the peer.
// The code argument to h is the received close code or CloseNoStatusReceived
// if the close message is empty. The default close handler sends a close
// message back to the peer.
//
// The handler function is called from the NextReader, ReadMessage and message
// reader Read methods. The application must read the connection to process
// close messages as described in the section on Control Messages above.
//
// The connection read methods return a CloseError when a close message is
// received. Most applications should handle close messages as part of their
// normal error handling. Applications should only set a close handler when the
// application must perform some action before sending a close message back to
// the peer.
func (c *WSConn) SetCloseHandler(h func(code int, text string) error) {
	c.c.SetCloseHandler(h)
}

// PingHandler returns the current ping handler
func (c *WSConn) PingHandler() func(appData string) error {
	return c.c.PingHandler()
}

// SetPingHandler sets the handler for ping messages received from the peer.
// The appData argument to h is the PING message application data. The default
// ping handler sends a pong to the peer.
//
// The handler function is called from the NextReader, ReadMessage and message
// reader Read methods. The application must read the connection to process
// ping messages as described in the section on Control Messages above.
func (c *WSConn) SetPingHandler(h func(appData string) error) {
	c.c.SetPingHandler(h)
}

// PongHandler returns the current pong handler
func (c *WSConn) PongHandler() func(appData string) error {
	return c.c.PongHandler()
}

// SetPongHandler sets the handler for pong messages received from the peer.
// The appData argument to h is the PONG message application data. The default
// pong handler does nothing.
//
// The handler function is called from the NextReader, ReadMessage and message
// reader Read methods. The application must read the connection to process
// pong messages as described in the section on Control Messages above.
func (c *WSConn) SetPongHandler(h func(appData string) error) {
	c.c.SetPongHandler(h)
}

// UnderlyingConn returns the internal net.Conn. This can be used to further
// modifications to connection specific flags.
func (c *WSConn) UnderlyingConn() net.Conn {
	return c.c.UnderlyingConn()
}

// EnableWriteCompression enables and disables write compression of
// subsequent text and binary messages. This function is a noop if
// compression was not negotiated with the peer.
func (c *WSConn) EnableWriteCompression(enable bool) {
	c.c.EnableWriteCompression(enable)
}

// SetCompressionLevel sets the flate compression level for subsequent text and
// binary messages. This function is a noop if compression was not negotiated
// with the peer. See the compress/flate package for a description of
// compression levels.
func (c *WSConn) SetCompressionLevel(level int) error {
	return c.c.SetCompressionLevel(level)
}

// WriteJSON writes the JSON encoding of v as a message.
//
// See the documentation for encoding/json Marshal for details about the
// conversion of Go values to JSON.
func (c *WSConn) WriteJSON(v interface{}) error {
	return c.c.WriteJSON(v)
}

// ReadJSON reads the next JSON-encoded message from the connection and stores
// it in the value pointed to by v.
//
// See the documentation for the encoding/json Unmarshal function for details
// about the conversion of JSON to a Go value.
func (c *WSConn) ReadJSON(v interface{}) error {
	return c.c.ReadJSON(v)
}

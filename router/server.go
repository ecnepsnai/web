package router

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/ecnepsnai/logtic"
)

// Server describes a server. Do not initialize a new copy of a Server{}, but instead use router.New()
type Server struct {
	impl       *impl
	httpServer *http.Server
	listener   *net.Listener
}

type impl struct {
	Lock                   *sync.RWMutex
	Index                  *endpoint
	NotFoundHandle         func(http.ResponseWriter, *http.Request)
	MethodNotAllowedHandle func(http.ResponseWriter, *http.Request)
	log                    *logtic.Source
}

// New will initialize a new Server instance and return it. This does not start the server.
func New() *Server {
	index := newEndpoint()
	log := logtic.Log.Connect("HTTP")
	s := &Server{
		impl: &impl{
			Lock:                   &sync.RWMutex{},
			Index:                  &index,
			NotFoundHandle:         defaultNotFoundHandle,
			MethodNotAllowedHandle: defaultMethodNotAllowedHandle,
			log:                    log,
		},
		httpServer: &http.Server{
			ReadTimeout:       5 * time.Minute,
			ReadHeaderTimeout: 5 * time.Minute,
			ErrorLog:          log.GoLogger(logtic.LevelError),
		},
	}
	return s
}

// ListenAndServe will listen for HTTP requests on the specified socket address.
// Valid addresses are typically in the form of: <IP Address>:<Port Number>. For IPv6 addresses, wrap the address in
// brackets.
//
// An error will only be returned if there was an error listening or the listener was closed.
func (s *Server) ListenAndServe(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		s.impl.log.PError("Error listening on address", map[string]interface{}{
			"address": addr,
			"error":   err.Error(),
		})
		return err
	}
	s.listener = &l
	s.httpServer.Handler = s.impl
	s.impl.log.PDebug("Listen", map[string]interface{}{
		"address": addr,
	})
	return s.Serve(l)
}

// Serve will listen for HTTP requests on the given listener.
//
// An error will only be returned if there was an error listening or the listener was abruptly closed.
func (s *Server) Serve(listener net.Listener) error {
	s.impl.log.Debug("Serve on listener")
	return http.Serve(listener, s.impl)
}

// Stop will stop the server. Server.ListenAndServe or Server.Serve will return net.ErrClosed. Does nothing if the
// was not listening or was already stopped.
func (s *Server) Stop() {
	if s.listener == nil {
		return
	}
	s.impl.log.Debug("Stopping server")
	l := *s.listener
	l.Close()
	s.impl.log.Info("Server stopped")
}

// SetNotFoundHandle will set the handle called when a request that did not match any registered path comes in.
//
// A default handle is set when the server is created.
func (s *Server) SetNotFoundHandle(handle func(w http.ResponseWriter, r *http.Request)) {
	s.impl.NotFoundHandle = handle
}

// SetMethodNotAllowedHandle will set the handle called when a request comes in for a known path but not the correct
// method.
//
// A default handle is set when the server is created.
func (s *Server) SetMethodNotAllowedHandle(handle func(w http.ResponseWriter, r *http.Request)) {
	s.impl.MethodNotAllowedHandle = handle
}

package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Server describes a API server instance
type Server struct {
	BindAddress  string
	router       *httprouter.Router
	socket       http.Server
	shuttingDown bool
}

// New create a new API server. Does not start the server.
func New(bindAddress string) *Server {
	httpRouter := httprouter.New()
	httpRouter.NotFound = notFoundHandler{}
	httpRouter.MethodNotAllowed = methodNotAllowedHandler{}
	return &Server{
		BindAddress: bindAddress,
		router:      httpRouter,
	}
}

// Start start the server. Blocks.
func (s *Server) Start() error {
	s.socket = http.Server{Addr: s.BindAddress, Handler: s.router}
	if err := s.socket.ListenAndServe(); err != nil {
		if s.shuttingDown {
			return nil
		}
		return err
	}
	return nil
}

// Stop stop the server.
func (s *Server) Stop() {
	s.shuttingDown = true
	s.socket.Close()
}

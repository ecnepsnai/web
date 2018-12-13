package web

import (
	"net/http"

	"github.com/ecnepsnai/logtic"

	"github.com/julienschmidt/httprouter"
)

// Server describes a API server instance
type Server struct {
	BindAddress  string
	router       *httprouter.Router
	socket       http.Server
	shuttingDown bool
	API          API
	HTTP         HTTP
	log          *logtic.Source
}

// New create a new API server. Does not start the server.
func New(bindAddress string) *Server {
	httpRouter := httprouter.New()
	server := Server{
		BindAddress: bindAddress,
		router:      httpRouter,
		log:         logtic.Connect("HTTP"),
	}
	httpRouter.NotFound = notFoundHandler{
		server: &server,
	}
	httpRouter.MethodNotAllowed = methodNotAllowedHandler{
		server: &server,
	}
	api := API{
		server: &server,
	}
	http := HTTP{
		server: &server,
	}
	server.API = api
	server.HTTP = http

	return &server
}

// Start start the server. Blocks.
func (s *Server) Start() error {
	s.socket = http.Server{Addr: s.BindAddress, Handler: s.router}
	s.log.Info("HTTP Server listening on %s", s.BindAddress)
	if err := s.socket.ListenAndServe(); err != nil {
		if s.shuttingDown {
			s.log.Info("HTTP server stopped")
			return nil
		}
		return err
	}
	return nil
}

// Stop stop the server.
func (s *Server) Stop() {
	s.log.Warn("Stopping HTTP server")
	s.shuttingDown = true
	s.socket.Close()
}

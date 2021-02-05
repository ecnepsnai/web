package web

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Server describes a API server instance
type Server struct {
	BindAddress             string
	router                  *httprouter.Router
	socket                  http.Server
	shuttingDown            bool
	API                     API
	HTTP                    HTTP
	NotFoundHandler         func(w http.ResponseWriter, r *http.Request)
	MethodNotAllowedHandler func(w http.ResponseWriter, r *http.Request)
}

// New create a new API server. Does not start the server.
func New(bindAddress string) *Server {
	httpRouter := httprouter.New()
	server := Server{
		BindAddress: bindAddress,
		router:      httpRouter,
	}
	httpRouter.NotFound = notFoundHandler{
		server: &server,
	}
	httpRouter.MethodNotAllowed = methodNotAllowedHandler{
		server: &server,
	}
	server.API = API{
		server: &server,
		limits: map[string]*rate.Limiter{},
	}
	server.HTTP = HTTP{
		server: &server,
	}

	return &server
}

// Start start the server. Blocks.
func (s *Server) Start() error {
	s.socket = http.Server{Addr: s.BindAddress, Handler: s.router}
	log.Info("HTTP Server listening on %s", s.BindAddress)
	if err := s.socket.ListenAndServe(); err != nil {
		if s.shuttingDown {
			log.Info("HTTP server stopped")
			return nil
		}
		return err
	}
	return nil
}

// Stop stop the server.
func (s *Server) Stop() {
	log.Warn("Stopping HTTP server")
	s.shuttingDown = true
	s.socket.Close()
}

type notFoundHandler struct {
	server *Server
}

func (n notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("HTTP %s %s -> %d", r.Method, r.RequestURI, 404)
	if n.server.NotFoundHandler != nil {
		n.server.NotFoundHandler(w, r)
		return
	}
	w.WriteHeader(404)
	w.Write([]byte("Not found"))
}

type methodNotAllowedHandler struct {
	server *Server
}

func (n methodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("HTTP %s %s -> %d", r.Method, r.RequestURI, 405)
	if n.server.MethodNotAllowedHandler != nil {
		n.server.MethodNotAllowedHandler(w, r)
		return
	}
	w.WriteHeader(405)
	w.Write([]byte("Method not allowed"))
}

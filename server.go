package web

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/ecnepsnai/logtic"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/time/rate"
)

// Server describes an web server
type Server struct {
	// The socket address that the server will listen to. Must be in the format of "address:port", such as
	// "localhost:8080" or "0.0.0.0:8080". Changing this after the server has started has no effect.
	BindAddress string
	// The port that this server is listening on.
	ListenPort uint16
	// The API instance that is used to register JSON endpoints.
	API API
	// The HTTP instance that is used to register plain HTTP endpoints.
	HTTP HTTP
	// The handler called when a request that does not match a registered path occurs. Defaults to a plain
	// HTTP 404 with "Not found" as the body.
	NotFoundHandler func(w http.ResponseWriter, r *http.Request)
	// The handler called when a request that did match a router but with the incorrect method occurs. Defaults to a
	// plain HTTP 405 with "Method not allowed" as the body.
	MethodNotAllowedHandler func(w http.ResponseWriter, r *http.Request)
	// The handler called when a request exceed the configured maximum per second limit. Defaults to a plain HTTP 429
	// with "Too many requests" as the body.
	RateLimitedHandler func(w http.ResponseWriter, r *http.Request)
	// Specify the maximum number of requests any given client IP address can make per second. Requests that are rate
	// limited will call the RateLimitedHandler, which you can override to customize the response.
	// Setting this to 0 disables rate limiting.
	MaxRequestsPerSecond int
	// The level to use when logging out HTTP requests. Maps to github.com/ecnepsnai/logtic levels. Defaults to Debug.
	RequestLogLevel int
	router          *httprouter.Router
	listener        net.Listener
	shuttingDown    bool
	limits          map[string]*rate.Limiter
	limitLock       *sync.Mutex
}

// New create a new server object that will bind to the provided address. Does not start the service automatically.
// Bind address must be in the format of "address:port", such as "localhost:8080" or "0.0.0.0:8080".
func New(bindAddress string) *Server {
	httpRouter := httprouter.New()
	server := Server{
		BindAddress:     bindAddress,
		RequestLogLevel: logtic.LevelDebug,
		router:          httpRouter,
		limits:          map[string]*rate.Limiter{},
		limitLock:       &sync.Mutex{},
	}
	httpRouter.NotFound = notFoundHandler{
		server: &server,
	}
	httpRouter.MethodNotAllowed = methodNotAllowedHandler{
		server: &server,
	}
	server.API = API{
		server: &server,
	}
	server.HTTP = HTTP{
		server: &server,
	}

	return &server
}

// Start will start the web server and listen on the socket address. This method blocks.
// If a server is stopped using the Stop() method, this returns no error.
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.BindAddress)
	if err != nil {
		log.PError("Error listening on address", map[string]interface{}{
			"listen_address": s.BindAddress,
			"error":          err.Error(),
		})
		return err
	}
	s.listener = listener
	s.ListenPort = uint16(listener.Addr().(*net.TCPAddr).Port)
	log.PInfo("HTTP server listen", map[string]interface{}{
		"listen_address": s.BindAddress,
		"listen_port":    s.ListenPort,
	})
	if err := http.Serve(listener, s.router); err != nil {
		if s.shuttingDown {
			log.Info("HTTP server stopped")
			return nil
		}
		return err
	}
	return nil
}

// Stop will stop the server. The Start() method will return without an error after stopping.
func (s *Server) Stop() {
	log.Warn("Stopping HTTP server")
	s.shuttingDown = true
	s.ListenPort = 0
	s.listener.Close()
}

type notFoundHandler struct {
	server *Server
}

func (n notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.PWrite(n.server.RequestLogLevel, "HTTP Request", map[string]interface{}{
		"remote_addr": getRealIP(r),
		"method":      r.Method,
		"url":         r.URL,
		"elapsed":     time.Duration(0).String(),
		"status":      404,
	})
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
	log.PWrite(n.server.RequestLogLevel, "HTTP Request", map[string]interface{}{
		"remote_addr": getRealIP(r),
		"method":      r.Method,
		"url":         r.URL,
		"elapsed":     time.Duration(0).String(),
		"status":      405,
	})
	if n.server.MethodNotAllowedHandler != nil {
		n.server.MethodNotAllowedHandler(w, r)
		return
	}
	w.WriteHeader(405)
	w.Write([]byte("Method not allowed"))
}

func (s *Server) isRateLimited(w http.ResponseWriter, r *http.Request) bool {
	// If rate limiting is not configured return a new limiter for each connection
	if s.MaxRequestsPerSecond == 0 {
		return false
	}

	s.limitLock.Lock()
	defer s.limitLock.Unlock()

	sourceIP := getRealIP(r).String()
	limiter := s.limits[sourceIP]
	if limiter == nil {
		// Allow MaxRequestsPerSecond every 1 second
		limiter = rate.NewLimiter(rate.Limit(s.MaxRequestsPerSecond), s.MaxRequestsPerSecond)
		s.limits[sourceIP] = limiter
	}

	if !limiter.Allow() {
		log.PWarn("Rate-limiting request", map[string]interface{}{
			"remote_addr": getRealIP(r),
			"method":      r.Method,
			"url":         r.URL,
		})
		log.PWrite(s.RequestLogLevel, "HTTP Request", map[string]interface{}{
			"remote_addr": getRealIP(r),
			"method":      r.Method,
			"url":         r.URL,
			"elapsed":     time.Duration(0).String(),
			"status":      429,
		})
		if s.RateLimitedHandler != nil {
			s.RateLimitedHandler(w, r)
		} else {
			w.WriteHeader(429)
			w.Write([]byte("Too many requests"))
		}
		return true
	}

	return false
}

package web

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/web/router"
	"golang.org/x/time/rate"
)

// Server describes an web server
type Server struct {
	// The socket address that the server is listening on. Only populated if the server was created with web.New().
	BindAddress string
	// The port that this server is listening on. Only populated if the server was created with web.New().
	ListenPort uint16
	// The JSON API server. API handles return data or an error, and all responses are wrapped in a common
	// response object; [web.JSONResponse].
	API API
	// HTTPEasy describes a easy HTTP server. HTTPEasy handles are expected to return a reader and specify the content
	// type and length themselves.
	//
	// The HTTPEasy server supports HTTP range requests, should the client request it and the application provide a
	// supported Reader [io.ReadSeekCloser].
	HTTPEasy HTTPEasy
	// The HTTP server. HTTP handles are exposed to the raw http request and response writers.
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
	// Additional options for the server
	Options ServerOptions

	router       *router.Server
	listener     net.Listener
	shuttingDown bool
	limits       map[string]*rate.Limiter
	limitLock    *sync.Mutex
}

type ServerOptions struct {
	// Specify the maximum number of requests any given client IP address can make per second. Requests that are rate
	// limited will call the RateLimitedHandler, which you can override to customize the response.
	// Setting this to 0 disables rate limiting.
	MaxRequestsPerSecond int
	// The level to use when logging out HTTP requests. Maps to github.com/ecnepsnai/logtic levels. Defaults to Debug.
	RequestLogLevel logtic.LogLevel
	// If true then the server will not try to reply with chunked data for a HTTP range request
	IgnoreHTTPRangeRequests bool
}

// New create a new server object that will bind to the provided address. Does not accept incoming connections until
// the server is started.
// Bind address must be in the format of "address:port", such as "localhost:8080" or "0.0.0.0:8080".
func New(bindAddress string) *Server {
	httpRouter := router.New()
	server := Server{
		BindAddress: bindAddress,
		Options: ServerOptions{
			RequestLogLevel: logtic.LevelDebug,
		},
		router:    httpRouter,
		limits:    map[string]*rate.Limiter{},
		limitLock: &sync.Mutex{},
	}
	httpRouter.SetNotFoundHandle(server.notFoundHandle)
	httpRouter.SetMethodNotAllowedHandle(server.methodNotAllowedHandle)
	server.API = API{
		server: &server,
	}
	server.HTTPEasy = HTTPEasy{
		server: &server,
	}
	server.HTTP = HTTP{
		server: &server,
	}

	return &server
}

// NewListener creates a new server object that will use the given listener. Does not accept incoming connections until
// the server is started.
func NewListener(listener net.Listener) *Server {
	httpRouter := router.New()
	server := Server{
		Options: ServerOptions{
			RequestLogLevel: logtic.LevelDebug,
		},
		router:    httpRouter,
		listener:  listener,
		limits:    map[string]*rate.Limiter{},
		limitLock: &sync.Mutex{},
	}
	httpRouter.SetNotFoundHandle(server.notFoundHandle)
	httpRouter.SetMethodNotAllowedHandle(server.methodNotAllowedHandle)
	server.API = API{
		server: &server,
	}
	server.HTTPEasy = HTTPEasy{
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
	if s.BindAddress != "" {
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
	}
	if err := s.router.Serve(s.listener); err != nil {
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

func (s *Server) notFoundHandle(w http.ResponseWriter, r *http.Request) {
	log.PWrite(s.Options.RequestLogLevel, "HTTP Request", map[string]interface{}{
		"remote_addr": RealRemoteAddr(r),
		"method":      r.Method,
		"url":         r.URL,
		"elapsed":     time.Duration(0).String(),
		"status":      404,
	})
	if s.NotFoundHandler != nil {
		s.NotFoundHandler(w, r)
		return
	}
	w.WriteHeader(404)
	w.Write([]byte("Not found"))
}

func (s *Server) methodNotAllowedHandle(w http.ResponseWriter, r *http.Request) {
	log.PWrite(s.Options.RequestLogLevel, "HTTP Request", map[string]interface{}{
		"remote_addr": RealRemoteAddr(r),
		"method":      r.Method,
		"url":         r.URL,
		"elapsed":     time.Duration(0).String(),
		"status":      405,
	})
	if s.MethodNotAllowedHandler != nil {
		s.MethodNotAllowedHandler(w, r)
		return
	}
	w.WriteHeader(405)
	w.Write([]byte("Method not allowed"))
}

func (s *Server) isRateLimited(w http.ResponseWriter, r *http.Request) bool {
	// If rate limiting is not configured return a new limiter for each connection
	if s.Options.MaxRequestsPerSecond == 0 {
		return false
	}

	s.limitLock.Lock()
	defer s.limitLock.Unlock()

	sourceIP := RealRemoteAddr(r).String()
	limiter := s.limits[sourceIP]
	if limiter == nil {
		// Allow MaxRequestsPerSecond every 1 second
		limiter = rate.NewLimiter(rate.Limit(s.Options.MaxRequestsPerSecond), s.Options.MaxRequestsPerSecond)
		s.limits[sourceIP] = limiter
	}

	if !limiter.Allow() {
		log.PWarn("Rate-limiting request", map[string]interface{}{
			"remote_addr": RealRemoteAddr(r),
			"method":      r.Method,
			"url":         r.URL,
		})
		log.PWrite(s.Options.RequestLogLevel, "HTTP Request", map[string]interface{}{
			"remote_addr": RealRemoteAddr(r),
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

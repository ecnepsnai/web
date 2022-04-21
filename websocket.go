package web

import (
	"encoding/json"
	"net/http"

	"github.com/ecnepsnai/web/router"
	"github.com/gorilla/websocket"
)

// Socket register a new websocket server at the given path
func (s *Server) Socket(path string, handle SocketHandle, options HandleOptions) {
	s.registerSocketEndpoint("GET", path, handle, options)
}

func (s *Server) registerSocketEndpoint(method string, path string, handle SocketHandle, options HandleOptions) {
	log.PDebug("Register websocket", map[string]interface{}{
		"method": method,
		"path":   path,
	})
	s.router.Handle(method, path, s.socketHandler(handle, options))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func socketPanicRecover() {
	if r := recover(); r != nil {
		log.Error("Recovered from socket handle panic: %#v", r)
	}
}

func (s *Server) socketHandler(endpointHandle SocketHandle, options HandleOptions) router.Handle {
	return func(w http.ResponseWriter, r router.Request) {
		defer socketPanicRecover()

		var userData interface{}

		if s.isRateLimited(w, r.HTTP) {
			return
		}

		if options.AuthenticateMethod != nil {
			userData = options.AuthenticateMethod(r.HTTP)
			if isUserdataNil(userData) {
				if options.UnauthorizedMethod == nil {
					log.PWarn("Rejected request to authenticated websocket endpoint", map[string]interface{}{
						"url":         r.HTTP.URL,
						"method":      r.HTTP.Method,
						"remote_addr": getRealIP(r.HTTP),
					})
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(Error{401, "Unauthorized"})
					return
				}

				options.UnauthorizedMethod(w, r.HTTP)
				return
			}
		}

		conn, err := upgrader.Upgrade(w, r.HTTP, nil)
		if err != nil {
			log.PError("Error upgrading client for websocket connection", map[string]interface{}{
				"error":       err.Error(),
				"remote_addr": getRealIP(r.HTTP),
			})
			return
		}
		endpointHandle(Request{
			Parameters: r.Parameters,
			UserData:   userData,
		}, WSConn{
			c: conn,
		})
		log.PWrite(s.Options.RequestLogLevel, "Websocket request", map[string]interface{}{
			"method":      r.HTTP.Method,
			"url":         r.HTTP.RequestURI,
			"remote_addr": getRealIP(r.HTTP),
		})
	}
}

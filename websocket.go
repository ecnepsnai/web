package web

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
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

func (s *Server) socketHandler(endpointHandle SocketHandle, options HandleOptions) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer socketPanicRecover()

		var userData interface{}

		if s.isRateLimited(w, r) {
			return
		}

		if options.AuthenticateMethod != nil {
			userData = options.AuthenticateMethod(r)
			if isUserdataNil(userData) {
				if options.UnauthorizedMethod == nil {
					log.PWarn("Rejected request to authenticated websocket endpoint", map[string]interface{}{
						"url":         r.URL,
						"method":      r.Method,
						"remote_addr": r.RemoteAddr,
					})
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(Error{401, "Unauthorized"})
					return
				}

				options.UnauthorizedMethod(w, r)
				return
			}
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.PError("Error upgrading client for websocket connection", map[string]interface{}{
				"error":       err.Error(),
				"remote_addr": r.RemoteAddr,
			})
			return
		}
		endpointHandle(Request{
			Params:   ps,
			UserData: userData,
		}, WSConn{
			c: conn,
		})
		log.PWrite(s.RequestLogLevel, "Websocket request", map[string]interface{}{
			"method":      r.Method,
			"url":         r.RequestURI,
			"remote_addr": r.RemoteAddr,
		})
	}
}

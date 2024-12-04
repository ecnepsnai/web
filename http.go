package web

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/ecnepsnai/web/router"
)

// HTTP describes a HTTP server. HTTP handles are exposed to the raw http request and response writers.
type HTTP struct {
	server *Server
}

// GET register a new HTTP GET request handle
func (h HTTP) GET(path string, handle HTTPHandle, options HandleOptions) {
	h.registerHTTPEndpoint("GET", path, handle, options)
}

// HEAD register a new HTTP HEAD request handle
func (h HTTP) HEAD(path string, handle HTTPHandle, options HandleOptions) {
	h.registerHTTPEndpoint("HEAD", path, handle, options)
}

// OPTIONS register a new HTTP OPTIONS request handle
func (h HTTP) OPTIONS(path string, handle HTTPHandle, options HandleOptions) {
	h.registerHTTPEndpoint("OPTIONS", path, handle, options)
}

// POST register a new HTTP POST request handle
func (h HTTP) POST(path string, handle HTTPHandle, options HandleOptions) {
	h.registerHTTPEndpoint("POST", path, handle, options)
}

// PUT register a new HTTP PUT request handle
func (h HTTP) PUT(path string, handle HTTPHandle, options HandleOptions) {
	h.registerHTTPEndpoint("PUT", path, handle, options)
}

// PATCH register a new HTTP PATCH request handle
func (h HTTP) PATCH(path string, handle HTTPHandle, options HandleOptions) {
	h.registerHTTPEndpoint("PATCH", path, handle, options)
}

// DELETE register a new HTTP DELETE request handle
func (h HTTP) DELETE(path string, handle HTTPHandle, options HandleOptions) {
	h.registerHTTPEndpoint("DELETE", path, handle, options)
}

func (h HTTP) registerHTTPEndpoint(method string, path string, handle HTTPHandle, options HandleOptions) {
	log.PDebug("Register HTTP endpoint", map[string]interface{}{
		"method": method,
		"path":   path,
	})
	h.server.router.Handle(method, path, h.httpPreHandle(handle, options))
}

func (h HTTP) httpPreHandle(endpointHandle HTTPHandle, options HandleOptions) router.Handle {
	return func(w http.ResponseWriter, request router.Request) {
		if options.PreHandle != nil {
			if err := options.PreHandle(w, request.HTTP); err != nil {
				return
			}
		}

		if h.server.isRateLimited(w, request.HTTP) {
			return
		}

		if options.MaxBodyLength > 0 {
			// We don't need to worry about this not being a number. Go's own HTTP server
			// won't respond to requests like these
			length, _ := strconv.ParseUint(request.HTTP.Header.Get("Content-Length"), 10, 64)

			if length > options.MaxBodyLength {
				log.PError("Rejecting HTTP request with oversized body", map[string]interface{}{
					"body_length": length,
					"max_length":  options.MaxBodyLength,
				})
				w.WriteHeader(413)
				return
			}
		}

		var userData interface{}
		if options.AuthenticateMethod != nil {
			userData = options.AuthenticateMethod(request.HTTP)
			if isUserdataNil(userData) {
				if options.UnauthorizedMethod == nil {
					log.PWarn("Rejected request to authenticated HTTP endpoint", map[string]interface{}{
						"url":         request.HTTP.URL,
						"method":      request.HTTP.Method,
						"remote_addr": RealRemoteAddr(request.HTTP),
					})
					w.Header().Set("Content-Type", "text/html")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("<html><head><title>Unauthorized</title></head><body><h1>Unauthorized</h1></body></html>"))
					return
				}

				options.UnauthorizedMethod(w, request.HTTP)
				return
			}
		}
		start := time.Now()
		defer func() {
			if p := recover(); p != nil {
				log.PError("Recovered from panic during HTTP handle", map[string]interface{}{
					"error":  fmt.Sprintf("%v", p),
					"route":  request.HTTP.URL.Path,
					"method": request.HTTP.Method,
					"stack":  string(debug.Stack()),
				})
				w.WriteHeader(500)
			}
		}()

		endpointHandle(w, Request{
			HTTP:       request.HTTP,
			Parameters: request.Parameters,
			UserData:   userData,
		})
		elapsed := time.Since(start)
		if !options.DontLogRequests {
			log.PWrite(h.server.Options.RequestLogLevel, "HTTP Request", map[string]interface{}{
				"remote_addr": RealRemoteAddr(request.HTTP),
				"method":      request.HTTP.Method,
				"url":         request.HTTP.URL,
				"elapsed":     elapsed.String(),
			})
		}
	}
}

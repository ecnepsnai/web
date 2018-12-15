package web

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// HTTP describes a HTTP server
type HTTP struct {
	server *Server
}

// HTTPHandle describes a method signature for handling an HTTP request
type HTTPHandle func(request Request) Response

// Static serve static files matching the request path to the given directory
func (h HTTP) Static(path string, directory string) {
	h.server.log.Debug("Serving files from '%s' matching path '%s'", directory, path)
	h.server.router.ServeFiles(path, http.Dir(directory))
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
	h.server.log.Debug("Register HTTP %s %s", method, path)
	h.server.router.Handle(method, path, h.httpAuthenticationHandler(handle, options))
}

func (h HTTP) httpAuthenticationHandler(endpointHandle HTTPHandle, options HandleOptions) httprouter.Handle {
	if options.AuthenticateMethod != nil {
		return func(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
			userData := options.AuthenticateMethod(request)
			if userData == nil {
				h.server.log.Warn("Rejected authenticated request")
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("<html><body><strong>Unauthorized</strong>"))
			} else {
				h.httpRequestHandler(endpointHandle, userData)(w, request, ps)
			}
		}
	}
	return h.httpRequestHandler(endpointHandle, nil)
}

func (h HTTP) httpRequestHandler(endpointHandle HTTPHandle, userData interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		request := Request{
			Writer:   w,
			HTTP:     r,
			Params:   ps,
			UserData: userData,
			log:      h.server.log,
		}
		start := time.Now()
		response := endpointHandle(request)
		elapsed := time.Since(start)

		if len(response.ContentType) == 0 {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		} else {
			w.Header().Set("Content-Type", response.ContentType)
		}

		for k, v := range response.Headers {
			w.Header().Set(k, v)
		}

		code := 200
		if response.Status != 0 {
			code = response.Status
		}
		h.server.log.Info("HTTP Request: %s %s -> %d (%s)", r.Method, r.RequestURI, code, elapsed)
		w.WriteHeader(code)

		if response.Reader != nil {
			readLength := 1024
			for {
				rbuf := make([]byte, readLength)
				length, err := response.Reader.Read(rbuf)
				if err != nil {
					if err.Error() == "EOF" {
						break
					}
					h.server.log.Error("Error reading response reader: %s", err.Error())
					w.WriteHeader(500)
					return
				}
				if length == 0 {
					break
				}
				_, err = w.Write(rbuf)
				if err != nil {
					h.server.log.Error("Error writing response: %s", err.Error())
					w.WriteHeader(500)
					return
				}
			}
			response.Reader.Close()
		}
	}
}

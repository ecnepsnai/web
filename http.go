package web

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

// HTTP describes a HTTP server
type HTTP struct {
	server *Server
}

// Static serve static files matching the request path to the given directory
func (h HTTP) Static(path string, directory string) {
	log.Debug("Serving files from '%s' matching path '%s'", directory, path)
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
	log.Debug("Register HTTP %s %s", method, path)
	h.server.router.Handle(method, path, h.httpPreHandle(handle, options))
}

func (h HTTP) httpPreHandle(endpointHandle HTTPHandle, options HandleOptions) httprouter.Handle {
	return func(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
		if options.MaxBodyLength > 0 {
			// We don't need to worry about this not being a number. Go's own HTTP server
			// won't respond to requests like these
			length, _ := strconv.ParseUint(request.Header.Get("Content-Length"), 10, 64)
			log.Debug("Body length: %d", length)

			if length > options.MaxBodyLength {
				log.Error("Rejecting HTTP request with body length %d", length)
				w.WriteHeader(413)
				return
			}
		}

		if options.AuthenticateMethod != nil {
			userData := options.AuthenticateMethod(request)
			if isUserdataNil(userData) {
				if options.UnauthorizedMethod == nil {
					log.Warn("Rejected authenticated request")
					w.Header().Set("Content-Type", "text/html")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("<html><head><title>Unauthorized</title></head><body><h1>Unauthorized</h1></body></html>"))
					return
				}

				options.UnauthorizedMethod(w, request)
			} else {
				h.httpPostHandle(endpointHandle, userData)(w, request, ps)
			}
			return
		}
		h.httpPostHandle(endpointHandle, nil)(w, request, ps)
	}
}

func (h HTTP) httpPostHandle(endpointHandle HTTPHandle, userData interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		request := Request{
			HTTP:     r,
			Params:   ps,
			UserData: userData,
			writer:   w,
		}
		start := time.Now()
		response := endpointHandle(request, Writer{w})
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
		log.Debug("HTTP Request: %s %s -> %d (%s)", r.Method, r.RequestURI, code, elapsed)
		w.WriteHeader(code)

		if response.Reader != nil {
			_, err := io.CopyBuffer(w, response.Reader, nil)
			response.Reader.Close()
			if err != nil {
				log.Error("Error writing response for HTTP request %s %s: %s", r.Method, r.RequestURI, err.Error())
				return
			}
		}
	}
}

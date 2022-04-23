package web

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ecnepsnai/web/router"
)

// HTTP describes a HTTP server. HTTP handles are expected to return a reader and specify the content
// type and length themselves.
//
// The HTTP server supports HTTP range requests, should the client request it and the application provide a
// supported Reader (io.ReadSeekCloser).
type HTTP struct {
	server *Server
}

// Static registers a GET and HEAD handle for all requests under path to serve any files matching the directory.
//
// For example:
//    directory = /usr/share/www/
//    path      = /static/
//
//    Request for '/static/image.jpg' would read file '/usr/share/www/image.jpg'
//
// Will panic if any handle is registered under path. Attempting to register a new handle under path after calling
// Static will panic.
//
// Caching will be enabled by default for all files served by this router. The mtime of the file will be used for the
// Last-Modified date.
//
// By default, the server will use the file extension (if any) to determine the MIME type for the response.
func (h HTTP) Static(path string, directory string) {
	log.PDebug("Serving files from directory", map[string]interface{}{
		"directory": directory,
		"path":      path,
	})
	h.server.router.ServeFiles(directory, path)
}

// GET register a new HTTP GET request handle
func (h HTTP) GET(path string, handle HTTPHandle, options HandleOptions) {
	h.registerHTTPEndpoint("GET", path, handle, options)
}

// HEAD register a new HTTP HEAD request handle
func (h HTTP) HEAD(path string, handle HTTPHandle, options HandleOptions) {
	h.registerHTTPEndpoint("HEAD", path, handle, options)
}

// GETHEAD registers both a HTTP GET and HTTP HEAD request handle. Equal to calling HTTP.GET and HTTP.HEAD.
//
// Handle responses can always return a reader, it will automatically be ignored for HEAD requests.
func (h HTTP) GETHEAD(path string, handle HTTPHandle, options HandleOptions) {
	h.registerHTTPEndpoint("GET", path, handle, options)
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

		if options.AuthenticateMethod != nil {
			userData := options.AuthenticateMethod(request.HTTP)
			if isUserdataNil(userData) {
				if options.UnauthorizedMethod == nil {
					log.PWarn("Rejected request to authenticated HTTP endpoint", map[string]interface{}{
						"url":         request.HTTP.URL,
						"method":      request.HTTP.Method,
						"remote_addr": getRealIP(request.HTTP),
					})
					w.Header().Set("Content-Type", "text/html")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("<html><head><title>Unauthorized</title></head><body><h1>Unauthorized</h1></body></html>"))
					return
				}

				options.UnauthorizedMethod(w, request.HTTP)
			} else {
				h.httpPostHandle(endpointHandle, userData)(w, request)
			}
			return
		}
		h.httpPostHandle(endpointHandle, nil)(w, request)
	}
}

func (h HTTP) httpPostHandle(endpointHandle HTTPHandle, userData interface{}) router.Handle {
	return func(w http.ResponseWriter, r router.Request) {
		request := Request{
			HTTP:       r.HTTP,
			Parameters: r.Parameters,
			UserData:   userData,
			writer:     w,
		}
		start := time.Now()
		response := endpointHandle(request, Writer{w})
		elapsed := time.Since(start)

		if response.Reader != nil {
			defer response.Reader.Close()
		}

		// Return a HTTP range response only if:
		// 1. A range was actually requested by the client
		// 2. The reader implemented Seek
		// 3. The response was either default or 200
		ranges := router.ParseRangeHeader(r.HTTP.Header.Get("range"))
		_, canSeek := response.Reader.(io.ReadSeekCloser)
		if len(ranges) > 0 && (response.Status == 0 || response.Status == 200) && !h.server.Options.IgnoreHTTPRangeRequests && canSeek {
			router.ServeHTTPRange(router.ServeHTTPRangeOptions{
				Headers:     response.Headers,
				Ranges:      ranges,
				Reader:      response.Reader.(io.ReadSeekCloser),
				TotalLength: response.ContentLength,
				MIMEType:    response.ContentType,
				Writer:      w,
			})
			log.PWrite(h.server.Options.RequestLogLevel, "HTTP Request", map[string]interface{}{
				"remote_addr": getRealIP(r.HTTP),
				"method":      r.HTTP.Method,
				"url":         r.HTTP.URL,
				"elapsed":     elapsed.String(),
				"status":      response.Status,
				"range":       true,
			})
			return
		}
		if canSeek && !h.server.Options.IgnoreHTTPRangeRequests {
			w.Header().Set("Accept-Ranges", "bytes")
		}

		if len(response.ContentType) == 0 {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		} else {
			w.Header().Set("Content-Type", response.ContentType)
		}

		if response.ContentLength > 0 {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", response.ContentLength))
		}

		for k, v := range response.Headers {
			w.Header().Set(k, v)
		}

		code := 200
		if response.Status != 0 {
			code = response.Status
		}
		log.PWrite(h.server.Options.RequestLogLevel, "HTTP Request", map[string]interface{}{
			"remote_addr": getRealIP(r.HTTP),
			"method":      r.HTTP.Method,
			"url":         r.HTTP.URL,
			"elapsed":     elapsed.String(),
			"status":      code,
		})
		w.WriteHeader(code)

		if r.HTTP.Method != "HEAD" && response.Reader != nil {
			_, err := io.CopyBuffer(w, response.Reader, nil)
			if err != nil {
				log.PError("Error writing response", map[string]interface{}{
					"method": r.HTTP.Method,
					"url":    r.HTTP.URL,
					"error":  err.Error(),
				})
				return
			}
		}
	}
}

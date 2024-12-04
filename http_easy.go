package web

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ecnepsnai/web/router"
)

// HTTPEasy describes a HTTPEasy server. HTTPEasy handles are expected to return a reader and specify the content
// type and length themselves.
//
// The HTTPEasy server supports HTTPEasy range requests, should the client request it and the application provide a
// supported Reader (io.ReadSeekCloser).
type HTTPEasy struct {
	server *Server
}

// Static registers a GET and HEAD handle for all requests under path to serve any files matching the directory.
//
// For example:
//
//	directory = /usr/share/www/
//	path      = /static/
//
//	Request for '/static/image.jpg' would read file '/usr/share/www/image.jpg'
//
// Will panic if any handle is registered under path. Attempting to register a new handle under path after calling
// Static will panic.
//
// Caching will be enabled by default for all files served by this router. The mtime of the file will be used for the
// Last-Modified date.
//
// By default, the server will use the file extension (if any) to determine the MIME type for the response.
func (h HTTPEasy) Static(path string, directory string) {
	log.PDebug("Serving files from directory", map[string]interface{}{
		"directory": directory,
		"path":      path,
	})
	h.server.router.ServeFiles(directory, path)
}

// GET register a new HTTP GET request handle
func (h HTTPEasy) GET(path string, handle HTTPEasyHandle, options HandleOptions) {
	h.registerHTTPEasyEndpoint("GET", path, handle, options)
}

// HEAD register a new HTTP HEAD request handle
func (h HTTPEasy) HEAD(path string, handle HTTPEasyHandle, options HandleOptions) {
	h.registerHTTPEasyEndpoint("HEAD", path, handle, options)
}

// GETHEAD registers both a HTTP GET and HTTP HEAD request handle. Equal to calling HTTPEasy.GET and HTTPEasy.HEAD.
//
// Handle responses can always return a reader, it will automatically be ignored for HEAD requests.
func (h HTTPEasy) GETHEAD(path string, handle HTTPEasyHandle, options HandleOptions) {
	h.registerHTTPEasyEndpoint("GET", path, handle, options)
	h.registerHTTPEasyEndpoint("HEAD", path, handle, options)
}

// OPTIONS register a new HTTP OPTIONS request handle
func (h HTTPEasy) OPTIONS(path string, handle HTTPEasyHandle, options HandleOptions) {
	h.registerHTTPEasyEndpoint("OPTIONS", path, handle, options)
}

// POST register a new HTTP POST request handle
func (h HTTPEasy) POST(path string, handle HTTPEasyHandle, options HandleOptions) {
	h.registerHTTPEasyEndpoint("POST", path, handle, options)
}

// PUT register a new HTTP PUT request handle
func (h HTTPEasy) PUT(path string, handle HTTPEasyHandle, options HandleOptions) {
	h.registerHTTPEasyEndpoint("PUT", path, handle, options)
}

// PATCH register a new HTTP PATCH request handle
func (h HTTPEasy) PATCH(path string, handle HTTPEasyHandle, options HandleOptions) {
	h.registerHTTPEasyEndpoint("PATCH", path, handle, options)
}

// DELETE register a new HTTP DELETE request handle
func (h HTTPEasy) DELETE(path string, handle HTTPEasyHandle, options HandleOptions) {
	h.registerHTTPEasyEndpoint("DELETE", path, handle, options)
}

func (h HTTPEasy) registerHTTPEasyEndpoint(method string, path string, handle HTTPEasyHandle, options HandleOptions) {
	log.PDebug("Register HTTP endpoint", map[string]interface{}{
		"method": method,
		"path":   path,
	})
	h.server.router.Handle(method, path, h.httpPreHandle(handle, options))
}

func (h HTTPEasy) httpPreHandle(endpointHandle HTTPEasyHandle, options HandleOptions) router.Handle {
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

		if options.AuthenticateMethod != nil {
			userData := options.AuthenticateMethod(request.HTTP)
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
			} else {
				h.httpPostHandle(endpointHandle, userData, options)(w, request)
			}
			return
		}
		h.httpPostHandle(endpointHandle, nil, options)(w, request)
	}
}

func (h HTTPEasy) httpPostHandle(endpointHandle HTTPEasyHandle, userData interface{}, options HandleOptions) router.Handle {
	return func(w http.ResponseWriter, r router.Request) {
		request := Request{
			HTTP:       r.HTTP,
			Parameters: r.Parameters,
			UserData:   userData,
		}
		start := time.Now()
		response := endpointHandle(request)
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
				"remote_addr": RealRemoteAddr(r.HTTP),
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

		for _, cookie := range response.Cookies {
			http.SetCookie(w, &cookie)
		}

		code := 200
		if response.Status != 0 {
			code = response.Status
		}
		if !options.DontLogRequests {
			log.PWrite(h.server.Options.RequestLogLevel, "HTTP Request", map[string]interface{}{
				"remote_addr": RealRemoteAddr(r.HTTP),
				"method":      r.HTTP.Method,
				"url":         r.HTTP.URL,
				"elapsed":     elapsed.String(),
				"status":      code,
			})
		}
		w.WriteHeader(code)

		if r.HTTP.Method != "HEAD" && response.Reader != nil {
			if copied, err := io.Copy(w, response.Reader); err != nil {
				if strings.Contains(err.Error(), "write: broken pipe") {
					return
				}

				log.PError("Error writing response data", map[string]interface{}{
					"method": r.HTTP.Method,
					"url":    r.HTTP.URL,
					"wrote":  copied,
					"error":  err.Error(),
				})
				return
			}
		}
	}
}

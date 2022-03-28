/*
Package router provides a simple & efficient parametrized HTTP router.

A HTTP router allows you to map a HTTP request method and path to a specific function. A parameterized HTTP router
allows you to designate specific portions of the request path as a parameter, which can later be fetched during the
request itself.

This package allows you modify the routing table ad-hoc, even while the server is running.
*/
package router

import (
	"net/http"
	"strings"
)

const (
	pathKeyIndex     = "__router_index"
	pathKeyParameter = "__router_parameter"
	pathKeyWildcard  = "__router_wildcard"
)

func init() {
	MimeGetter = &extensionMimeGetterType{}
}

// Request describes a HTTP request
type Request struct {
	// The underlaying HTTP request
	HTTP *http.Request
	// A map of any parameters from the router path mapped to their values from the request path
	Parameters map[string]string
}

// Handle describes the signature for a handle of a path
type Handle func(http.ResponseWriter, Request)

type endpoint struct {
	Methods   map[string]Handle
	Children  map[string]endpoint
	Parameter string
}

func newEndpoint() endpoint {
	return endpoint{
		Methods:  map[string]Handle{},
		Children: map[string]endpoint{},
	}
}

func (s *impl) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.Lock.RLock()
	defer func() {
		s.Lock.RUnlock()
		if r := recover(); r != nil {
			s.log.PError("Recovered from router panic", map[string]interface{}{
				"request_method": req.Method,
				"request_path":   req.URL.Path,
				"error":          r.(string),
			})
			w.WriteHeader(500)
		}
	}()

	// Handle wildcard roots
	if wildcardChild, exists := s.Index.Children[pathKeyWildcard]; exists {
		handler, present := wildcardChild.Methods[req.Method]
		if !present {
			s.MethodNotAllowedHandle(w, req)
			return
		}
		handler(w, Request{req, map[string]string{
			wildcardChild.Parameter: req.URL.Path[1:], // trim the leading /
		}})
		return
	}

	parameters := map[string]string{}

	// If the request path ends in a slash, append the index path key
	path := req.URL.Path
	if path[len(path)-1] == '/' {
		path += pathKeyIndex
	}
	segments := strings.Split(path[1:], "/")

	// Start at the index:
	// - if there is a matching child for the segment fetch that
	// - if there is no matching segment, check for a wildcard
	// - if there is no wildcard, check for a parameter
	// once we've reached the last segment, find the handle for the request method
	// - if no method found, check for other methods to return a 405
	parent := s.Index
	for i, segment := range segments {
		child, exists := parent.Children[segment]

		if !exists {
			if wildcardChild, exists := parent.Children[pathKeyWildcard]; exists {
				handler, present := wildcardChild.Methods[req.Method]
				if !present {
					s.MethodNotAllowedHandle(w, req)
					return
				}
				value := strings.Join(segments[i:], "/")
				if req.URL.Path[len(req.URL.Path)-1] == '/' {
					value = value[0 : len(value)-len(pathKeyIndex)]
				}
				parameters[wildcardChild.Parameter] = value
				handler(w, Request{req, parameters})
				return
			}
			parameterChild, exists := parent.Children[pathKeyParameter]
			if !exists {
				s.NotFoundHandle(w, req)
				return
			}
			child = parameterChild
			parameters[parameterChild.Parameter] = segment
		}

		parent = &child

		if i == len(segments)-1 { // last segment
			handler, present := parent.Methods[req.Method]
			if !present {
				if len(parent.Methods) > 0 {
					s.MethodNotAllowedHandle(w, req)
					return
				}

				s.NotFoundHandle(w, req)
				return
			}

			handler(w, Request{req, parameters})
			return
		}
	}

	// should never actually hit this
	s.NotFoundHandle(w, req)
}

func (s *Server) registerHandle(method, path string, handler Handle) {
	s.impl.Lock.Lock()
	defer s.impl.Lock.Unlock()

	if path[len(path)-1] == '/' {
		path += pathKeyIndex
	}
	segments := strings.Split(path[1:], "/")

	parent := s.impl.Index
	for i, segment := range segments {
		parameter := ""

		// Since you can only have one unique parameter per segment, we don't
		// have to worry about what the parameter name is.
		if len(segment) > 1 {
			if segment[0] == '*' {
				parameter = segment[1:]
				segment = pathKeyWildcard
				i = len(segments) - 1
			} else if segment[0] == ':' {
				parameter = segment[1:]
				segment = pathKeyParameter
			}
		}

		if wc, exists := parent.Children[pathKeyWildcard]; exists {
			if wc.Parameter != parameter {
				panic("Path segment collides with wildcard")
			}
		}

		child, exists := parent.Children[segment]
		if !exists {
			if segment == pathKeyWildcard && len(parent.Children) >= 1 {
				panic("Path segment collides with wildcard")
			}
			if len(parent.Children) == 1 && parent.Children[pathKeyParameter].Parameter != "" {
				panic("Path part '" + segment + "' collides with existing parameter :" + child.Parameter)
			}

			child = newEndpoint()
			child.Parameter = parameter
			parent.Children[segment] = child
		}

		parent = &child

		if i == len(segments)-1 {
			if _, exists := parent.Methods[method]; exists {
				panic("Handle already registered for method and path")
			}
			parent.Methods[method] = handler
			s.impl.log.PDebug("Register handle", map[string]interface{}{
				"method": method,
				"path":   path,
			})
			return
		}
	}
}

// Handle registers a handler for an HTTP request of method to path.
//
// Method must be a valid HTTP method, in all caps. Path must always begin with a forward slash /. Will panic on invalid
// vales. Will panic if registering a duplicate method & path.
//
// Handle may be called even while the server is listening and is threadsafe.
//
// Any segment that begins with a colon (:) will be parameterized. The value of all parameters for the path will
// be populated into the Parameters map included in the Request object in the handler. For example:
//
//     handle path  = "/widgets/:widget_id/cost/:currency"
//     request path = "/widgets/1234/cost/cad"
//     parameters   = { "widget_id": "1234", "currency": "cad" }
//
// Any segment that begins with an astreisk (*) will be parameterized as well, however unlike colon parameters, these
// will include the entire remaining path as the value of the parameter, whereas colon parameters will only include that
// segment as the value. Multiple methods can be registered for the same wildcard path, provided they use the same
// parameter name. Any segments included after the parameter name are ignored.
// For example
//
//     handle path  = "/proxy/*url"
//     request path = "/proxy/some/multi/segmented/value"
//     parameters   = { "url": "some/multi/segmented/value" }
//
// Parameter segments are exclusive, meaning you can not have a static segment at the same position as a
// parameterized element. For example, these both will panic:
//
//     // This panics because /all occupies the same segment as the parameter :username
//     server.Handle("GET", "/users/:username", ...)
//     server.Handle("GET", "/users/all", ...)
//
//     // This panics because /user/id occupied the same segment as the wildcard parameter *param
//     server.Handle("GET", "/users/*param", ...)
//     server.Handle("GET", "/users/user/id", ...)
//
// Paths that end with a slash are unique to those that don't. For example, these would be considred unique by the
// router:
//
//     server.Handle("GET", "/users/all/", ...)
//     server.Handle("GET", "/users/all", ...)
//
func (s *Server) Handle(method, path string, handler Handle) {
	methods := map[string]bool{
		"CONNECT": true,
		"DELETE":  true,
		"GET":     true,
		"HEAD":    true,
		"OPTIONS": true,
		"PATCH":   true,
		"POST":    true,
		"PUT":     true,
		"TRACE":   true,
	}
	if _, m := methods[method]; !m {
		panic("Invalid HTTP method " + method)
	}
	if path[0] != '/' {
		panic("Path must start with /")
	}
	if strings.Contains(path, pathKeyIndex) || strings.Contains(path, pathKeyParameter) || strings.Contains(path, pathKeyWildcard) {
		panic("Path contains reserved string sequence")
	}

	s.registerHandle(method, path, handler)
}

// RemoveHandle will remove any handler for the given method and path. If no handle exists, it does nothing.
// If both method and path are * it removes everything from the routing table.
//
// Note that parameter names are not considered when removing a path. For example, you may register a path with
// `/:username` and remove it with `/:something_else`.
//
// This may be called even while the server is listening and is threadsafe.
func (s *Server) RemoveHandle(method, path string) {
	s.impl.Lock.Lock()
	defer s.impl.Lock.Unlock()

	if method == "*" && path == "*" {
		s.impl.log.Debug("Removing all handles")
		index := newEndpoint()
		s.impl.Index = &index
		return
	}

	if path == "" || path[0] != '/' {
		return
	}

	if path[len(path)-1] == '/' {
		path += pathKeyIndex
	}
	segments := strings.Split(path[1:], "/")

	parent := s.impl.Index
	for i, segment := range segments {
		if len(segment) > 1 {
			if segment[0] == '*' {
				segment = pathKeyWildcard
			} else if segment[0] == ':' {
				segment = pathKeyParameter
			}
		}

		child, exists := parent.Children[segment]
		if !exists {
			return
		}

		if i == len(segments)-1 {
			delete(child.Methods, method)
			s.impl.log.PDebug("Remove handle", map[string]interface{}{
				"method": method,
				"path":   path,
			})
			if len(child.Methods) == 0 {
				delete(parent.Children, segment)
			}
			return
		} else {
			parent = &child
		}
	}
}

// ServeFiles registers a handler for all requests under urlRoot to serve any files matching the same path in
// a local filesystem directory localRoot.
//
// For example:
//    localRoot = /usr/share/www/
//    urlRoot   = /static/
//
//    Request for '/static/image.jpg' would read file '/usr/share/www/image.jpg'
//
// Will panic if any handle is registered under urlRoot. Attempting to register a new handle under urlRoot after calling
// ServeFiles will panic.
//
// Caching will be enabled by default for all files served by this router. The mtime of the file will be used for the
// Last-Modified date.
//
// By default, the server will use the file extension (if any) to determine the MIME type for the response. You may
// use your own MIME detection by implementing the IMime interface and setting MimeGetter.
//
// The server will also instruct clients to cache files served for up-to 1 day. You can control this with the
// CacheMaxAge variable.
//
// When a directory is requested, the router will look for an index file with the name from the IndexFileName variable.
// If no file is found, a directory listing will automatically be generated. You can control this with the
// GenerateDirectoryListing variable.
func (s *Server) ServeFiles(localRoot string, urlRoot string) {
	var handle Handle = func(rw http.ResponseWriter, r Request) {
		s.impl.serveStatic(localRoot, r.Parameters["path"], rw, r.HTTP)
	}

	if urlRoot[len(urlRoot)-1] != '/' {
		urlRoot += "/"
	}
	urlRoot += "*path"

	s.Handle("GET", urlRoot, handle)
	s.Handle("HEAD", urlRoot, handle)
}

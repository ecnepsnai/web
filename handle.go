package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Handle describes a method signature for handling an API request
type Handle func(request Request) (interface{}, *Error)

// HandleOptions describes options for a route
type HandleOptions struct {
	AuthenticateMethod func(request *http.Request) interface{}
}

// APIGET register a new HTTP GET request handle
func (s *Server) APIGET(path string, handle Handle, options HandleOptions) {
	s.registerAPIEndpoint("GET", path, handle, options)
}

// APIHEAD register a new HTTP HEAD request handle
func (s *Server) APIHEAD(path string, handle Handle, options HandleOptions) {
	s.registerAPIEndpoint("HEAD", path, handle, options)
}

// APIOPTIONS register a new HTTP OPTIONS request handle
func (s *Server) APIOPTIONS(path string, handle Handle, options HandleOptions) {
	s.registerAPIEndpoint("OPTIONS", path, handle, options)
}

// APIPOST register a new HTTP POST request handle
func (s *Server) APIPOST(path string, handle Handle, options HandleOptions) {
	s.registerAPIEndpoint("POST", path, handle, options)
}

// APIPUT register a new HTTP PUT request handle
func (s *Server) APIPUT(path string, handle Handle, options HandleOptions) {
	s.registerAPIEndpoint("PUT", path, handle, options)
}

// APIPATCH register a new HTTP PATCH request handle
func (s *Server) APIPATCH(path string, handle Handle, options HandleOptions) {
	s.registerAPIEndpoint("PATCH", path, handle, options)
}

// APIDELETE register a new HTTP DELETE request handle
func (s *Server) APIDELETE(path string, handle Handle, options HandleOptions) {
	s.registerAPIEndpoint("DELETE", path, handle, options)
}

func (s *Server) registerAPIEndpoint(method string, path string, handle Handle, options HandleOptions) {
	s.router.Handle(method, path, s.apiAuthenticationHandler(handle, options))
}

func (s *Server) apiAuthenticationHandler(endpointHandle Handle, options HandleOptions) httprouter.Handle {
	if options.AuthenticateMethod != nil {
		return func(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
			userData := options.AuthenticateMethod(request)
			if userData == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(Error{401, "Unauthorized"})
			} else {
				s.apiRequestHandler(endpointHandle, userData)(w, request, ps)
			}
		}
	}
	return s.apiRequestHandler(endpointHandle, nil)
}

func (s *Server) apiRequestHandler(endpointHandle Handle, userData interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		response := Response{}
		request := Request{
			Writer:   w,
			HTTP:     r,
			Params:   ps,
			UserData: userData,
		}

		data, err := endpointHandle(request)
		if err != nil {
			response.Code = err.Code
			w.WriteHeader(err.Code)
			response.Error = *err
		} else {
			response.Code = 200
			response.Data = data
		}
		json.NewEncoder(w).Encode(response)
	}
}

type notFoundHandler struct{}

func (n notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Error{404, "Not found"})
}

type methodNotAllowedHandler struct{}

func (n methodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(405)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Error{405, "Method not allowed"})
}

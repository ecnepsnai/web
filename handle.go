package web

import (
	"net/http"
	"reflect"
)

// APIHandle describes a method signature for handling an API request
type APIHandle func(request Request) (interface{}, *Error)

// HTTPHandle describes a method signature for handling an HTTP request
type HTTPHandle func(request Request, writer Writer) Response

// SocketHandle describes a method signature for handling a HTTP websocket request
type SocketHandle func(request Request, conn WSConn)

// HandleOptions describes options for a route
type HandleOptions struct {
	// AuthenticateMethod method called to determine if a request is properly authenticated or not.
	// Optional - Omit this entirely if no authentication is needed for the request.
	// Return nil to signal an unauthenticated request, which will be rejected.
	// Objects returned will be passed to the handle as the UserData object.
	AuthenticateMethod func(request *http.Request) interface{}
	// UnauthorizedMethod method called when an unauthenticated request occurs (AuthenticateMethod returned nil)
	// to customize the response seen by the user.
	// Optional - Omit this to have a default response.
	UnauthorizedMethod func(w http.ResponseWriter, request *http.Request)
	// MaxBodyLength defines the maximum length accepted for any HTTP request body. Requests that
	// exceed this limit will receive a 413 Payload Too Large response.
	// The default value of 0 will not reject requests with large bodies.
	MaxBodyLength uint64
}

func isUserdataNil(userData interface{}) bool {
	return userData == nil || (reflect.ValueOf(userData).Kind() == reflect.Ptr && reflect.ValueOf(userData).IsNil())
}

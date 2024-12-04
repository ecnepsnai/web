package web

import (
	"net/http"
	"reflect"
)

// APIHandle describes a method signature for handling an API request
type APIHandle func(request Request) (interface{}, *APIResponse, *Error)

// HTTPEasyHandle describes a method signature for handling an HTTP request
type HTTPEasyHandle func(request Request) HTTPResponse

// HTTPHandle describes a method signature for handling an HTTP request
type HTTPHandle func(w http.ResponseWriter, r Request)

// SocketHandle describes a method signature for handling a HTTP websocket request
type SocketHandle func(request Request, conn *WSConn)

// HandleOptions describes options for a route
type HandleOptions struct {
	// AuthenticateMethod method called to determine if a request is properly authenticated or not. If a method is
	// provided, then it is called for each incoming request. The value returned by this method is passed as the
	// UserData field of a [web.Request]. Returning nil signals an unauthenticated request, which will be handled by
	// the UnauthorizedMethod (if provided) or a default handle. If the AuthenticateMethod is not provided, then the
	// UserData field is nil.
	AuthenticateMethod func(request *http.Request) interface{}
	// PreHandle is an optional method that is called immediately upon receiving the HTTP request, before authentication
	// and before rate limit checks. This method allows servers to provide early handling of a request before any
	// processing happens.
	//
	// The returned error is only used as a nil check, the value of any error isn't used. If an error is returned then
	// no more processing is performed. It is assumed that a response will have been written to w.
	//
	// If nil is returned then the request will continue normally, no status should have been written to w. Any headers
	// added may be overwritten by the handle.
	PreHandle func(w http.ResponseWriter, request *http.Request) error
	// UnauthorizedMethod method called when an unauthenticated request occurs, i.e.AuthenticateMethod returned nil,
	// which allows you to customize the response seen by the user.
	// If omitted, a default handle is used.
	UnauthorizedMethod func(w http.ResponseWriter, request *http.Request)
	// MaxBodyLength defines the maximum length accepted for any HTTP request body. Requests that exceed this limit will
	// receive a "413 Payload Too Large" response. The default value of 0 will not reject requests with large bodies.
	MaxBodyLength uint64
	// DontLogRequests if true then requests to this handle are not logged
	DontLogRequests bool
}

func isUserdataNil(userData interface{}) bool {
	return userData == nil || (reflect.ValueOf(userData).Kind() == reflect.Ptr && reflect.ValueOf(userData).IsNil())
}

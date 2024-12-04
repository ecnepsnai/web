package web

import (
	"encoding/json"
	"net"
	"net/http"
)

// Request describes an API request
type Request struct {
	// The original HTTP request
	HTTP *http.Request
	// URL path parameters (not query parameters). Keys do not include the ':' or '*'.
	Parameters map[string]string
	// User data provided from the result of the AuthenticateRequest method on the handle options
	UserData any
}

// Decoder describes a generic interface that has a Decode function
type Decoder interface {
	Decode(v any) error
}

// DecodeJSON unmarshal the JSON body to the provided interface
func (r Request) DecodeJSON(v any) *Error {
	return r.Decode(v, json.NewDecoder(r.HTTP.Body))
}

// Decode will unmarshal the request body to v using the given decoder
func (r Request) Decode(v any, decoder Decoder) *Error {
	if err := json.NewDecoder(r.HTTP.Body).Decode(v); err != nil {
		log.PError("Invalid request", map[string]interface{}{
			"error": err.Error(),
		})
		return CommonErrors.BadRequest
	}

	return nil
}

// RealRemoteAddr will try to get the real IP address of the incoming connection taking proxies into
// consideration. This function looks for the `X-Real-IP`, `X-Forwarded-For`, and `CF-Connecting-IP`
// headers, and if those don't exist will return the remote address of the connection.
func (r Request) RealRemoteAddr() net.IP {
	return RealRemoteAddr(r.HTTP)
}

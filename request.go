package web

import (
	"encoding/json"
	"net"
	"net/http"
)

// Request describes an API request
type Request struct {
	HTTP       *http.Request
	Parameters map[string]string
	UserData   interface{}
}

// DecodeJSON unmarshal the JSON body to the provided interface
func (r Request) DecodeJSON(v interface{}) *Error {
	if err := json.NewDecoder(r.HTTP.Body).Decode(v); err != nil {
		log.PError("Invalid JSON request", map[string]interface{}{
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

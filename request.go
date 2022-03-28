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
	writer     http.ResponseWriter
}

// AddCookie add a cookie to the response
func (r Request) AddCookie(cookie *http.Cookie) {
	http.SetCookie(r.writer, cookie)
}

// Writer HTTP response writer
type Writer struct {
	w http.ResponseWriter
}

// Write see http.ResponseWriter.Write for more
func (w Writer) Write(d []byte) (int, error) {
	return w.w.Write(d)
}

// WriteHeader see http.ResponseWriter.Write for more
func (w Writer) WriteHeader(statusCode int) {
	w.w.WriteHeader(statusCode)
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

// ClientIPAddress returns the IP address of the client. It supports the 'X-Real-IP' and 'X-Forwarded-For' headers.
func (r Request) ClientIPAddress() net.IP {
	return getRealIP(r.HTTP)
}

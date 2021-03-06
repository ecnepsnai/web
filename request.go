package web

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Request describes an API request
type Request struct {
	HTTP     *http.Request
	Params   httprouter.Params
	UserData interface{}
	writer   http.ResponseWriter
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

// Decode unmarshal the JSON body to the provided interface
func (r Request) Decode(v interface{}) *Error {
	if err := json.NewDecoder(r.HTTP.Body).Decode(v); err != nil {
		log.Error("Invalid JSON request: error='%s'", err.Error())
		return CommonErrors.BadRequest
	}

	return nil
}

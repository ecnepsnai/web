package web

import (
	"encoding/json"
	"net/http"

	"github.com/ecnepsnai/logtic"
	"github.com/julienschmidt/httprouter"
)

// Request describes an API request
type Request struct {
	HTTP     *http.Request
	Params   httprouter.Params
	UserData interface{}
	log      *logtic.Source
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
		r.log.Error("Invalid JSON request: %s", err.Error())
		return CommonErrors.BadRequest
	}

	return nil
}

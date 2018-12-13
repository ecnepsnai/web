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
	Writer   http.ResponseWriter
	UserData interface{}
}

// Decode unmarshal the JSON body to the provided interface
func (r Request) Decode(v interface{}) *Error {
	if err := json.NewDecoder(r.HTTP.Body).Decode(v); err != nil {
		return CommonErrors.BadRequest
	}

	return nil
}

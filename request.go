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
	Writer   http.ResponseWriter
	UserData interface{}
	log      *logtic.Source
}

// Decode unmarshal the JSON body to the provided interface
func (r Request) Decode(v interface{}) *Error {
	if err := json.NewDecoder(r.HTTP.Body).Decode(v); err != nil {
		r.log.Error("Invalid JSON request: %s", err.Error())
		return CommonErrors.BadRequest
	}

	return nil
}

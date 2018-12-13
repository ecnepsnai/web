package web

import "io"

// JSONResponse describes an API response object
type JSONResponse struct {
	Code  int         `json:"code,omitempty"`
	Error Error       `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

// Response describes a HTTP response
type Response struct {
	Reader      io.ReadCloser
	Status      int
	Headers     map[string]string
	ContentType string
}

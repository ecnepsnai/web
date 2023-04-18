package web

import (
	"io"
	"net/http"
)

// APIResponse describes additional response properties for API handles
type APIResponse struct {
	// Additional headers to append to the response.
	Headers map[string]string
	// Cookies to set on the response.
	Cookies []http.Cookie
}

// JSONResponse describes an API response object
type JSONResponse struct {
	// The actual data of the response
	Data interface{} `json:"data,omitempty"`
	// If an error occured, details about the error
	Error *Error `json:"error,omitempty"`
	// The HTTP status code for the response
	//
	// Deprecated: will be removed in the next breaking update
	Code int `json:"code"`
}

// HTTPResponse describes a HTTP response
type HTTPResponse struct {
	// The reader for the response. Will be closed when the HTTP response is finished. Can be nil.
	//
	// If a io.ReadSeekCloser is provided then ranged data may be provided for a HTTP range request.
	Reader io.ReadCloser
	// The status code for the response. If 0 then 200 is implied.
	Status int
	// Additional headers to append to the response.
	Headers map[string]string
	// Cookies to set on the response.
	Cookies []http.Cookie
	// The content type of the response. Will overwrite any 'content-type' header in Headers.
	ContentType string
	// The length of the content.
	ContentLength uint64
}

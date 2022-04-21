package web

import "io"

// JSONResponse describes an API response object
type JSONResponse struct {
	// The actual data of the response
	Data interface{} `json:"data,omitempty"`
	// If an error occured, details about the error
	Error *Error `json:"error,omitempty"`
	// The HTTP status code for the response
	Code int `json:"code"`
}

// HTTPResponse describes a HTTP response
type HTTPResponse struct {
	// Reader the reader for the response. Will be closed when the HTTP response is finished. Can be nil.
	//
	// If a io.ReadSeekCloser is provided then ranged data may be provided for a HTTP range request.
	Reader io.ReadCloser
	// Status the status code for the response. If 0 then 200 is implied.
	Status int
	// Headers any additional headers to append to the response.
	Headers map[string]string
	// ContentType the content type of the response. Will overwrite any 'content-type' header in Headers.
	ContentType string
	// ContentLength the length of the content. If zero then Status should be 204.
	ContentLength uint64
}

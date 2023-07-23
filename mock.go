package web

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// Parameters for creating a mock request for uses in tests
type MockRequestParameters struct {
	// User data to be passed into the handler. May be nil.
	UserData interface{}
	// URL parameters (not query parameters) to be populated into the request.Params object in the handler. May be nil.
	Parameters map[string]string
	// Object to be encoded with JSON as the body. May be nil. Exclusive to Body.
	JSONBody interface{}
	// Body data. May be nil. Exclusive to JSONBody.
	Body io.ReadCloser
	// Optional HTTP request to pass to the handler.
	Request *http.Request
}

// MockRequest will generate a mock request for testing your handlers. Will panic for invalid parameters.
func MockRequest(parameters MockRequestParameters) Request {
	var httpRequest *http.Request

	if parameters.Request != nil {
		httpRequest = parameters.Request
	} else {
		httpRequest = &http.Request{}
	}

	httpRequest.RemoteAddr = "[::1]:65535"

	if parameters.JSONBody != nil && parameters.Body != nil {
		panic("cannot provide both JSON and data body")
	}

	if parameters.JSONBody != nil {
		b := &bytes.Buffer{}
		if err := json.NewEncoder(b).Encode(parameters.JSONBody); err != nil {
			panic(err)
		}
		httpRequest.Body = io.NopCloser(b)
	}

	if parameters.Body != nil {
		httpRequest.Body = parameters.Body
	}

	return Request{
		HTTP:       httpRequest,
		Parameters: parameters.Parameters,
		UserData:   parameters.UserData,
	}
}

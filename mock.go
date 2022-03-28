package web

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// MockRequest generate a mock request for testing your handlers. Body will be encoded as JSON and may panic if invalid.
func MockRequest(userData interface{}, params map[string]string, body interface{}) Request {
	var data []byte
	if body != nil {
		d, err := json.Marshal(&body)
		if err != nil {
			panic(err)
		}
		data = d
	}

	r := Request{
		HTTP: &http.Request{
			RemoteAddr: "[::1]:65535",
			Body:       io.NopCloser(bytes.NewReader(data)),
			Header: http.Header{
				"User-Agent": []string{"go test"},
			},
		},
		Parameters: params,
	}

	if userData != nil {
		r.UserData = userData
	}

	return r
}

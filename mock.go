package web

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// MockRequest generate a mock request for testing your handlers. Body will be encoded as JSON and may panic if invalid.
func MockRequest(userData interface{}, params map[string]string, body interface{}) Request {
	var p []httprouter.Param

	for k, v := range params {
		p = append(p, httprouter.Param{
			Key:   k,
			Value: v,
		})
	}

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
			Body: io.NopCloser(bytes.NewReader(data)),
			Header: http.Header{
				"User-Agent": []string{"go test"},
			},
		},
		Params: p,
	}

	if userData != nil {
		r.UserData = userData
	}

	return r
}

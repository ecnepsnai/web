package web_test

import (
	"testing"

	"github.com/ecnepsnai/web"
)

func TestMock(t *testing.T) {
	type exampleType struct {
		Enabled bool
	}

	userData := 1

	handle := func(request web.Request) (interface{}, *web.Error) {
		example := exampleType{}

		if err := request.DecodeJSON(&example); err != nil {
			t.Error("Error decoding example JSON object from mocked request")
		}
		if !example.Enabled {
			t.Error("Invalid HTTP body from mocked request")
		}

		if request.UserData.(int) != userData {
			t.Error("Invalid user data")
		}

		if request.Parameters["foo"] != "bar" {
			t.Error("Invalid request path parameters")
		}

		return nil, nil
	}

	request := web.MockRequest(web.MockRequestParameters{
		UserData:   userData,
		Parameters: map[string]string{"foo": "bar"},
		JSONBody:   exampleType{true},
	})
	handle(request)
}

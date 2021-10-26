package web_test

import "github.com/ecnepsnai/web"

func ExampleValidationError() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request) (interface{}, *web.Error) {
		username := request.Parameters["username"]

		return nil, web.ValidationError("No user with username %s", username)
	}
	server.API.GET("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

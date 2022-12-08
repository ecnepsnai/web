package web_test

import "github.com/ecnepsnai/web"

func ExampleAPI_GET() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		username := request.Parameters["username"]

		return map[string]string{
			"username": username,
		}, nil, nil
	}
	server.API.GET("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

func ExampleAPI_HEAD() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		return nil, nil, nil
	}
	server.API.HEAD("/users/user/", handle, web.HandleOptions{})

	server.Start()
}

func ExampleAPI_OPTIONS() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		return nil, nil, nil
	}
	server.API.OPTIONS("/users/user/", handle, web.HandleOptions{})

	server.Start()
}

func ExampleAPI_POST() {
	server := web.New("127.0.0.1:8080")

	type userRequestType struct {
		FirstName string `json:"first_name"`
	}

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		username := request.Parameters["username"]
		params := userRequestType{}
		if err := request.DecodeJSON(&params); err != nil {
			return nil, nil, err
		}

		return map[string]string{
			"first_name": params.FirstName,
			"username":   username,
		}, nil, nil
	}
	server.API.POST("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

func ExampleAPI_PUT() {
	server := web.New("127.0.0.1:8080")

	type userRequestType struct {
		FirstName string `json:"first_name"`
	}

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		username := request.Parameters["username"]
		params := userRequestType{}
		if err := request.DecodeJSON(&params); err != nil {
			return nil, nil, err
		}

		return map[string]string{
			"first_name": params.FirstName,
			"username":   username,
		}, nil, nil
	}
	server.API.PUT("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

func ExampleAPI_PATCH() {
	server := web.New("127.0.0.1:8080")

	type userRequestType struct {
		FirstName string `json:"first_name"`
	}

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		username := request.Parameters["username"]
		params := userRequestType{}
		if err := request.DecodeJSON(&params); err != nil {
			return nil, nil, err
		}

		return map[string]string{
			"first_name": params.FirstName,
			"username":   username,
		}, nil, nil
	}
	server.API.PATCH("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

func ExampleAPI_DELETE() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		username := request.Parameters["username"]

		return map[string]string{
			"username": username,
		}, nil, nil
	}
	server.API.DELETE("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

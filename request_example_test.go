package web_test

import (
	"fmt"

	"github.com/ecnepsnai/web"
)

func ExampleRequest_DecodeJSON() {
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

func ExampleRequest_RealRemoteAddr() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		clientAddr := request.RealRemoteAddr().String()
		fmt.Printf("%s\n", clientAddr)
		return clientAddr, nil, nil
	}
	server.API.POST("/ip/my_ip", handle, web.HandleOptions{})

	server.Start()
}

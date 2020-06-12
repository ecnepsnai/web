package web_test

import "github.com/ecnepsnai/web"

func ExampleRequest_Decode() {
	server := web.New("127.0.0.1:8080")

	type userRequestType struct {
		FirstName string `json:"first_name"`
	}

	handle := func(request web.Request) (interface{}, *web.Error) {
		username := request.Params.ByName("username")
		params := userRequestType{}
		if err := request.Decode(&params); err != nil {
			return nil, err
		}

		return map[string]string{
			"first_name": params.FirstName,
			"username":   username,
		}, nil
	}
	server.API.POST("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

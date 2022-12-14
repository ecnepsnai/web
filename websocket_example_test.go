package web_test

import "github.com/ecnepsnai/web"

func ExampleServer_Socket() {
	server := web.New("127.0.0.1:8080")

	type questionType struct {
		Name string
	}

	type answerType struct {
		Reply string
	}

	handle := func(request web.Request, conn *web.WSConn) {
		question := questionType{}
		if err := conn.ReadJSON(&question); err != nil {
			return
		}

		reply := answerType{
			Reply: "Hello, " + question.Name,
		}
		if err := conn.WriteJSON(&reply); err != nil {
			return
		}
	}

	options := web.HandleOptions{}
	server.Socket("/greeting", handle, options)

	server.Start()
}

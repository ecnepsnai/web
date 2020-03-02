/*
Package web is a simple HTTP server in Golang that is designed for both front and back-end web applications.

It includes a powerful JSON-based REST API framework and a simple-to-use HTTP router
for serving non-JSON content (like HTML).

It includes simple controls to allow for user authentication with contextual data
being avaialble in every request.

JSON API Example

	server = web.New("127.0.0.1:8080")
	if err := server.Start(); err != nil {
		panic(err)
	}

	handle := func(request web.Request) (interface{}, *Error) {
		return time.Now.Unix(), nil
	}
	options := web.HandleOptions{}
	server.API.GET("/time", handle, options)

File-Serving Example

	server = web.New("127.0.0.1:8080")
	if err := server.Start(); err != nil {
		panic(err)
	}

	handle := func(request web.Request, writer web.Writer) web.Response {
		f, err := os.Open("/foo/bar")
		if err != nil {
			return CommonErrors.ServerError
		}
		return Response{
			Reader: f,
		}
	}
	options := HandleOptions{}
	server.HTTP.GET("/file", handle, options)

Authentication Example

	server = New("127.0.0.1:8080")
	if err := server.Start(); err != nil {
		panic(err)
	}

	userInfo := User{
		Username: "ecnepsnai",
	}

	// Login
	loginHandle := func(request Request) (interface{}, *Error) {
		cookie := http.Cookie{
			Name:    "session",
			Value:   "1",
			Path:    "/",
			Expires: time.Now().AddDate(0, 0, 1),
		}
		http.SetCookie(request.Writer, &cookie)
		return true, nil
	}
	unauthenticatedOptions := HandleOptions{}
	server.API.GET("/login", loginHandle, unauthenticatedOptions)

	// Get User Info
	userHandle := func(request Request) (interface{}, *Error) {
		user := request.UserData.(User)
		return user, nil
	}

	authenticatedOptions := HandleOptions{
		// The authenticate method is where you can pass contextual data to the request
		// return nil to indicate authentication failure
		AuthenticateMethod: func(request *http.Request) interface{} {
			cookie, err := request.Cookie("session")
			if err != nil || cookie == nil {
				return nil
			}
			if cookie.Value != "1" {
				return nil
			}
			return map[string]string{
				"foo": "bar",
			}
		},
	}
	server.API.GET("/user", userHandle, authenticatedOptions)

Websocket Example

	server = web.New("127.0.0.1:8080")
	if err := server.Start(); err != nil {
		panic(err)
	}

	type questionType struct{
		Name string
	}

	type answerType struct{
		Reply string
	}

	handle := func(request web.Request, conn web.WSConn) {
		question := questionType{}
		if err := conn.ReadJSON(&question); err != nil {
			return
		}

		reply := answerType{
			Reply: "Hello, " + question.Name
		}
		if err := conn.WriteJSON(&reply); err != nil {
			return
		}
	}
	options := web.HandleOptions{}
	server.Socket("/greeting", handle, options)
*/
package web

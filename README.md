# Web

[![Go Report Card](https://goreportcard.com/badge/github.com/ecnepsnai/web?style=flat-square)](https://goreportcard.com/report/github.com/ecnepsnai/web)
[![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/ecnepsnai/web)
[![Releases](https://img.shields.io/github/release/ecnepsnai/web/all.svg?style=flat-square)](https://github.com/ecnepsnai/web/releases)
[![LICENSE](https://img.shields.io/github/license/ecnepsnai/web.svg?style=flat-square)](https://github.com/ecnepsnai/web/blob/master/LICENSE)

Web is a simple HTTP server in Golang that is designed for both front and back-end web
applications.

It includes a powerful JSON-based REST API framework, a simple-to-use HTTP router
for serving non-JSON content (like HTML), and an interface for Websockets.

It includes simple controls to allow for user authentication with contextual data
being avaialble in every request.

## JSON API Example

Return a JSON object with the UNIX Epoch when users browse to `/time`

```golang
server = web.New("127.0.0.1:8080")
if err := server.Start(); err != nil {
	panic(err)
}

handle := func(request web.Request) (interface{}, *Error) {
	return time.Now.Unix(), nil
}
options := web.HandleOptions{}
server.API.GET("/time", handle, options)
```

## File-Serving Example

Return the contents of the file `/foo/bar` when users browse to `/file`

```golang
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
```

## Authentication Example

```golang
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
```

## Websocket Example

A JSON-Based websocket server that replies to users questions

```golang
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
```
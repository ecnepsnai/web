package web_test

import (
	"net/http"
	"os"
	"time"

	"github.com/ecnepsnai/web"
)

func Example_json() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request) (interface{}, *web.Error) {
		return time.Now().Unix(), nil
	}

	options := web.HandleOptions{}
	server.API.GET("/time", handle, options)

	if err := server.Start(); err != nil {
		panic(err)
	}
}

func Example_file() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request, writer web.Writer) web.Response {
		f, err := os.Open("/foo/bar")
		if err != nil {
			return web.Response{
				Status: 500,
			}
		}
		return web.Response{
			Reader: f,
		}
	}

	options := web.HandleOptions{}
	server.HTTP.GET("/file", handle, options)

	if err := server.Start(); err != nil {
		panic(err)
	}
}

func Example_authentication() {
	server := web.New("127.0.0.1:8080")

	type User struct {
		Username string `json:"username"`
	}

	// Login
	loginHandle := func(request web.Request) (interface{}, *web.Error) {
		// Do any authentication logic here

		// Assuming the user authenticated successfully...
		request.AddCookie(&http.Cookie{
			Name:    "session",
			Value:   "1",
			Path:    "/",
			Expires: time.Now().AddDate(0, 0, 1),
		})
		return true, nil
	}
	unauthenticatedOptions := web.HandleOptions{}
	server.API.GET("/login", loginHandle, unauthenticatedOptions)

	// Get User Info
	getUserHandle := func(request web.Request) (interface{}, *web.Error) {
		user := request.UserData.(User)
		return user, nil
	}

	authenticatedOptions := web.HandleOptions{
		// The authenticate method is where you validate that a request if from an authenticated, or simple "logged in"
		// user. In this example, we validate that a cookie is present.
		// Any data returned by this method is provided into the request handler as Request.UserData
		// Returning nil results in a HTTP 403 response
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
	// Notice that we used a different HandleOptions instance with our AuthenticateMethod
	// an options without any AuthenticateMethod is considered unauthenticated
	server.API.GET("/user", getUserHandle, authenticatedOptions)

	if err := server.Start(); err != nil {
		panic(err)
	}
}

func Example_websocket() {
	server := web.New("127.0.0.1:8080")

	type questionType struct {
		Name string
	}

	type answerType struct {
		Reply string
	}

	handle := func(request web.Request, conn web.WSConn) {
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

	if err := server.Start(); err != nil {
		panic(err)
	}
}

func Example_ratelimit() {
	server := web.New("127.0.0.1:8080")

	// Restrict each connecting IP address to a maximum of 5 requests per second
	server.MaxRequestsPerSecond = 5

	// Handle called when a request is rejected due to rate limiting
	server.RateLimitedHandler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(429)
		w.Write([]byte("Too many requests"))
	}

	handle := func(request web.Request) (interface{}, *web.Error) {
		return time.Now().Unix(), nil
	}

	options := web.HandleOptions{}
	server.API.GET("/time", handle, options)

	if err := server.Start(); err != nil {
		panic(err)
	}
}

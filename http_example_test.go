package web_test

import (
	"os"

	"github.com/ecnepsnai/web"
)

func ExampleHTTP_Static() {
	server := web.New("127.0.0.1:8080")

	server.HTTP.Static("/static/*", "/path/to/static/files")

	server.Start()
}

func ExampleHTTP_GET() {
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
	server.HTTP.GET("/users/user", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTP_HEAD() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request, writer web.Writer) web.Response {
		return web.Response{
			Headers: map[string]string{
				"X-Fancy-Header": "some value",
			},
		}
	}
	server.HTTP.HEAD("/users/user", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTP_OPTIONS() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request, writer web.Writer) web.Response {
		return web.Response{
			Headers: map[string]string{
				"X-Fancy-Header": "some value",
			},
		}
	}
	server.HTTP.OPTIONS("/users/user", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTP_POST() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request, writer web.Writer) web.Response {
		username := request.Params.ByName("username")

		f, err := os.Open("/foo/bar")
		if err != nil {
			return web.Response{
				Status: 500,
			}
		}
		return web.Response{
			Headers: map[string]string{
				"X-Username": username,
			},
			Reader: f,
		}
	}
	server.HTTP.POST("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTP_PUT() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request, writer web.Writer) web.Response {
		username := request.Params.ByName("username")

		f, err := os.Open("/foo/bar")
		if err != nil {
			return web.Response{
				Status: 500,
			}
		}
		return web.Response{
			Headers: map[string]string{
				"X-Username": username,
			},
			Reader: f,
		}
	}
	server.HTTP.PUT("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTP_PATCH() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request, writer web.Writer) web.Response {
		username := request.Params.ByName("username")

		f, err := os.Open("/foo/bar")
		if err != nil {
			return web.Response{
				Status: 500,
			}
		}
		return web.Response{
			Headers: map[string]string{
				"X-Username": username,
			},
			Reader: f,
		}
	}
	server.HTTP.PATCH("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTP_DELETE() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request, writer web.Writer) web.Response {
		username := request.Params.ByName("username")
		return web.Response{
			Headers: map[string]string{
				"X-Username": username,
			},
		}
	}
	server.HTTP.DELETE("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

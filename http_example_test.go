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

	handle := func(request web.Request, writer web.Writer) web.HTTPResponse {
		f, err := os.Open("/foo/bar")
		if err != nil {
			return web.HTTPResponse{
				Status: 500,
			}
		}
		return web.HTTPResponse{
			Reader: f,
		}
	}
	server.HTTP.GET("/users/user", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTP_HEAD() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request, writer web.Writer) web.HTTPResponse {
		return web.HTTPResponse{
			Headers: map[string]string{
				"X-Fancy-Header": "some value",
			},
		}
	}
	server.HTTP.HEAD("/users/user", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTP_GETHEAD() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request, writer web.Writer) web.HTTPResponse {
		f, err := os.Open("/foo/bar")
		info, ierr := f.Stat()
		if err != nil || ierr != nil {
			return web.HTTPResponse{
				Status: 500,
			}
		}
		return web.HTTPResponse{
			Reader:        f, // the file will not be read for HTTP HEAD requests, but it will be closed.
			ContentType:   "text/plain",
			ContentLength: uint64(info.Size()),
		}
	}
	server.HTTP.GETHEAD("/users/user", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTP_OPTIONS() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request, writer web.Writer) web.HTTPResponse {
		return web.HTTPResponse{
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

	handle := func(request web.Request, writer web.Writer) web.HTTPResponse {
		username := request.Parameters["username"]

		f, err := os.Open("/foo/bar")
		if err != nil {
			return web.HTTPResponse{
				Status: 500,
			}
		}
		return web.HTTPResponse{
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

	handle := func(request web.Request, writer web.Writer) web.HTTPResponse {
		username := request.Parameters["username"]

		f, err := os.Open("/foo/bar")
		if err != nil {
			return web.HTTPResponse{
				Status: 500,
			}
		}
		return web.HTTPResponse{
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

	handle := func(request web.Request, writer web.Writer) web.HTTPResponse {
		username := request.Parameters["username"]

		f, err := os.Open("/foo/bar")
		if err != nil {
			return web.HTTPResponse{
				Status: 500,
			}
		}
		return web.HTTPResponse{
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

	handle := func(request web.Request, writer web.Writer) web.HTTPResponse {
		username := request.Parameters["username"]
		return web.HTTPResponse{
			Headers: map[string]string{
				"X-Username": username,
			},
		}
	}
	server.HTTP.DELETE("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

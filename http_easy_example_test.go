package web_test

import (
	"os"

	"github.com/ecnepsnai/web"
)

func ExampleHTTPEasy_Static() {
	server := web.New("127.0.0.1:8080")

	server.HTTPEasy.Static("/static/*", "/path/to/static/files")

	server.Start()
}

func ExampleHTTPEasy_GET() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request) web.HTTPResponse {
		f, err := os.Open("/foo/bar")
		info, ierr := f.Stat()
		if err != nil || ierr != nil {
			return web.HTTPResponse{
				Status: 500,
			}
		}
		return web.HTTPResponse{
			Reader:        f, // The file will be closed automatically
			ContentType:   "text/plain",
			ContentLength: uint64(info.Size()),
		}
	}
	server.HTTPEasy.GET("/users/user", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTPEasy_HEAD() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request) web.HTTPResponse {
		return web.HTTPResponse{
			Headers: map[string]string{
				"X-Fancy-Header": "some value",
			},
		}
	}
	server.HTTPEasy.HEAD("/users/user", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTPEasy_GETHEAD() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request) web.HTTPResponse {
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
	server.HTTPEasy.GETHEAD("/users/user", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTPEasy_OPTIONS() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request) web.HTTPResponse {
		return web.HTTPResponse{
			Headers: map[string]string{
				"X-Fancy-Header": "some value",
			},
		}
	}
	server.HTTPEasy.OPTIONS("/users/user", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTPEasy_POST() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request) web.HTTPResponse {
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
	server.HTTPEasy.POST("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTPEasy_PUT() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request) web.HTTPResponse {
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
	server.HTTPEasy.PUT("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTPEasy_PATCH() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request) web.HTTPResponse {
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
	server.HTTPEasy.PATCH("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTPEasy_DELETE() {
	server := web.New("127.0.0.1:8080")

	handle := func(request web.Request) web.HTTPResponse {
		username := request.Parameters["username"]
		return web.HTTPResponse{
			Headers: map[string]string{
				"X-Username": username,
			},
		}
	}
	server.HTTPEasy.DELETE("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

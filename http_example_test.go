package web_test

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ecnepsnai/web"
)

func ExampleHTTP_GET() {
	server := web.New("127.0.0.1:8080")

	handle := func(w http.ResponseWriter, r web.Request) {
		f, _ := os.Open("/foo/bar")
		info, _ := f.Stat()
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))
		io.Copy(w, f)
	}
	server.HTTP.GET("/users/user", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTP_HEAD() {
	server := web.New("127.0.0.1:8080")

	handle := func(w http.ResponseWriter, r web.Request) {
		w.Header().Set("X-Fancy-Header", "Some value")
		w.WriteHeader(204)
	}
	server.HTTP.HEAD("/users/user", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTP_OPTIONS() {
	server := web.New("127.0.0.1:8080")

	handle := func(w http.ResponseWriter, r web.Request) {
		w.Header().Set("X-Fancy-Header", "Some value")
		w.WriteHeader(200)
	}
	server.HTTP.OPTIONS("/users/user", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTP_POST() {
	server := web.New("127.0.0.1:8080")

	handle := func(w http.ResponseWriter, r web.Request) {
		username := r.Parameters["username"]

		w.Header().Set("X-Username", username)
		w.WriteHeader(200)
	}
	server.HTTP.POST("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTP_PUT() {
	server := web.New("127.0.0.1:8080")

	handle := func(w http.ResponseWriter, r web.Request) {
		username := r.Parameters["username"]

		w.Header().Set("X-Username", username)
		w.WriteHeader(200)
	}
	server.HTTP.PUT("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTP_PATCH() {
	server := web.New("127.0.0.1:8080")

	handle := func(w http.ResponseWriter, r web.Request) {
		username := r.Parameters["username"]

		w.Header().Set("X-Username", username)
		w.WriteHeader(200)
	}
	server.HTTP.PATCH("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

func ExampleHTTP_DELETE() {
	server := web.New("127.0.0.1:8080")

	handle := func(w http.ResponseWriter, r web.Request) {
		username := r.Parameters["username"]

		w.Header().Set("X-Username", username)
		w.WriteHeader(200)
	}
	server.HTTP.DELETE("/users/user/:username", handle, web.HandleOptions{})

	server.Start()
}

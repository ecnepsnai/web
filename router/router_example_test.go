package router_test

import (
	"net"
	"net/http"

	"github.com/ecnepsnai/web/router"
)

func ExampleServer_Handle() {
	server := router.New()
	server.Handle("GET", "/hello/:greeting", func(rw http.ResponseWriter, r router.Request) {
		rw.Write([]byte("Hello, " + r.Parameters["greeting"]))
	})
}

func ExampleServer_RemoveHandle() {
	server := router.New()
	server.Handle("GET", "/hello/:greeting", func(rw http.ResponseWriter, r router.Request) {
		rw.Write([]byte("Hello, " + r.Parameters["greeting"]))
	})
	server.RemoveHandle("GET", "/hello/:greeting")
	server.RemoveHandle("*", "*") // Will remove everything from the routing table!
}

func ExampleServer_ListenAndServe() {
	server := router.New()
	server.ListenAndServe("127.0.0.1:8080")
}

func ExampleServer_Serve() {
	server := router.New()
	l, _ := net.Listen("tcp", "[::1]:8080")
	server.Serve(l)
}

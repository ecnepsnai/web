package router_test

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/ecnepsnai/web/router"
)

func TestDirectoryIndex(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(path.Join(dir, "example.txt"), []byte("foo"), os.ModePerm)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/example/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	testStaticRequest(t, "GET", "http://"+listenAddress+"/example/", 200, "text/html; charset=utf-8")
}

func TestDirectoryIndexFile(t *testing.T) {
	router.IndexFileName = "index.htm"

	dir := t.TempDir()
	os.WriteFile(path.Join(dir, router.IndexFileName), []byte("foo"), os.ModePerm)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/example/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	testStaticRequest(t, "GET", "http://"+listenAddress+"/example/"+router.IndexFileName, 200, "text/html")

	router.IndexFileName = "index.html"
}

func TestDirectoryIndexDisabled(t *testing.T) {
	router.GenerateDirectoryListing = false

	dir := t.TempDir()
	os.WriteFile(path.Join(dir, "example.txt"), []byte("foo"), os.ModePerm)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/example/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	testStaticRequest(t, "GET", "http://"+listenAddress+"/example/", 404, "text/plain; charset=utf-8")

	router.GenerateDirectoryListing = true
}

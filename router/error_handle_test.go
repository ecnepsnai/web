package router_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/ecnepsnai/web/router"
)

func testURLContentType(t *testing.T, method, url, accept, expectedContentType string, expectedStatusCode int) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Accept", accept)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Unexpected status code for URL '%s'. Expected %d got %d", url, expectedStatusCode, resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != expectedContentType {
		t.Errorf("Unexpected content type for URL '%s'. Expected '%s' got '%s'", url, expectedContentType, resp.Header.Get("Content-Type"))
	}
}

func TestErrorHandleNotFound(t *testing.T) {
	t.Parallel()

	listenAddress := getListenAddress()

	server := router.New()
	server.Handle("GET", "/", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	testURLContentType(t, "GET", "http://"+listenAddress+"/", "", "", 200)
	testURLContentType(t, "GET", "http://"+listenAddress+"/foo", "text/html", "text/html; charset=utf-8", 404)
}

func TestErrorHandleMethodNotAllowed(t *testing.T) {
	t.Parallel()

	listenAddress := getListenAddress()

	server := router.New()
	server.Handle("GET", "/", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	testURLContentType(t, "GET", "http://"+listenAddress+"/", "", "", 200)
	testURLContentType(t, "POST", "http://"+listenAddress+"/", "text/html", "text/html; charset=utf-8", 405)
}

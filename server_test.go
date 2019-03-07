package web

import (
	"net/http"
	"testing"
)

func TestRestartServer(t *testing.T) {
	handle := func(request Request) (interface{}, *Error) {
		return true, nil
	}
	options := HandleOptions{}

	path1 := randomString(5)
	path2 := randomString(5)

	server.API.GET("/"+path1, handle, options)

	check := func(path string, expected int) {
		resp, err := http.Get("http://localhost:9557/" + path)
		if err != nil {
			t.Errorf("Network error: %s", err.Error())
		}

		if resp.StatusCode != expected {
			t.Errorf("Unexpected status code. Expected %d got %d", expected, resp.StatusCode)
		}
	}

	check(path1, 200)
	check(path2, 404)

	server.Stop()
	server.API.GET("/"+path2, handle, options)
	testStartServer()

	check(path1, 200)
	check(path2, 200)
}

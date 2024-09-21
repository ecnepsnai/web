package web_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ecnepsnai/web"
)

func TestRequestRealIP(t *testing.T) {
	t.Parallel()
	server := newServer()

	var expectedIP string

	handle := func(request web.Request) web.HTTPResponse {
		if request.RealRemoteAddr().String() != expectedIP {
			t.Errorf("Unexpected client IP address. Expected '%s' got '%s'", expectedIP, request.RealRemoteAddr().String())
		}
		return web.HTTPResponse{}
	}
	options := web.HandleOptions{}

	path := randomString(5)
	server.HTTPEasy.GET("/"+path, handle, options)

	var req *http.Request
	var err error

	// X-Real-IP
	expectedIP = "1.1.1.1"
	req, err = http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path), nil)
	if err != nil {
		t.Fatalf("Error forming request: %s", err.Error())
	}
	req.Header.Add("X-Real-IP", expectedIP)
	if _, err := http.DefaultClient.Do(req); err != nil {
		t.Fatalf("Network error: %s", err.Error())
	}

	// X-Forwarded-For
	expectedIP = "1::1"
	req, err = http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path), nil)
	if err != nil {
		t.Fatalf("Error forming request: %s", err.Error())
	}
	req.Header.Add("X-Forwarded-For", expectedIP)
	if _, err := http.DefaultClient.Do(req); err != nil {
		t.Fatalf("Network error: %s", err.Error())
	}

	// Client IP
	expectedIP = "::1"
	req, err = http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path), nil)
	if err != nil {
		t.Fatalf("Error forming request: %s", err.Error())
	}
	if _, err := http.DefaultClient.Do(req); err != nil {
		t.Fatalf("Network error: %s", err.Error())
	}
}

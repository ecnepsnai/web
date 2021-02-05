package web_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ecnepsnai/web"
)

func TestRestartServer(t *testing.T) {
	t.Parallel()
	server := newServer()

	handle := func(request web.Request) (interface{}, *web.Error) {
		return true, nil
	}
	options := web.HandleOptions{}

	path1 := randomString(5)
	path2 := randomString(5)

	server.API.GET("/"+path1, handle, options)

	check := func(path string, expected int) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path))
		if err != nil {
			t.Fatalf("Network error: %s", err.Error())
		}

		if resp.StatusCode != expected {
			t.Fatalf("Unexpected status code. Expected %d got %d", expected, resp.StatusCode)
		}
	}

	check(path1, 200)
	check(path2, 404)

	server.Stop()

	// Check it's actually stopped
	if _, err := http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path1)); err == nil {
		t.Fatalf("No error returned when one expected")
	}

	server.API.GET("/"+path2, handle, options)
	go server.Start()
	i := 0
	for i < 10 {
		if server.ListenPort > 0 {
			break
		}
		i++
		time.Sleep(5 * time.Millisecond)
	}
	if server.ListenPort == 0 {
		panic("Server didn't start in time")
	}

	check(path1, 200)
	check(path2, 200)

	server.Stop()
}

func TestNotFoundHandle(t *testing.T) {
	t.Parallel()
	server := newServer()

	htmlResponse := "<html><body><p>Not found</p></body></html>"
	jsonResponse := "{\"error\": \"not found\"}"
	plainResponse := "not found"

	server.NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		accept := r.Header.Get("Accept")
		if strings.Contains(accept, "text/html") {
			w.Write([]byte(htmlResponse))
		} else if strings.Contains(accept, "application/json") {
			w.Write([]byte(jsonResponse))
		} else {
			w.Write([]byte(plainResponse))
		}
	}

	// HTML
	func() {
		req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, randomString(6)), nil)
		req.Header.Add("Accept", "text/html")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Network error: %s", err.Error())
		}
		if resp.StatusCode != 404 {
			t.Fatalf("Unexpected status code. Expected %d Got %d", 404, resp.StatusCode)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != htmlResponse {
			t.Fatalf("Unexpected body %v", body)
		}
	}()

	// JSON
	func() {
		req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, randomString(6)), nil)
		req.Header.Add("Accept", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Network error: %s", err.Error())
		}
		if resp.StatusCode != 404 {
			t.Fatalf("Unexpected status code. Expected %d Got %d", 404, resp.StatusCode)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != jsonResponse {
			t.Fatalf("Unexpected body %v", body)
		}
	}()

	// Plain
	func() {
		req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, randomString(6)), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Network error: %s", err.Error())
		}
		if resp.StatusCode != 404 {
			t.Fatalf("Unexpected status code. Expected %d Got %d", 404, resp.StatusCode)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != plainResponse {
			t.Fatalf("Unexpected body %v", body)
		}
	}()
}

func TestMethodNotAllowed(t *testing.T) {
	t.Parallel()
	server := newServer()

	htmlResponse := "<html><body><p>method not allowed</p></body></html>"
	jsonResponse := "{\"error\": \"method not allowed\"}"
	plainResponse := "method not allowed"

	server.MethodNotAllowedHandler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		accept := r.Header.Get("Accept")
		if strings.Contains(accept, "text/html") {
			w.Write([]byte(htmlResponse))
		} else if strings.Contains(accept, "application/json") {
			w.Write([]byte(jsonResponse))
		} else {
			w.Write([]byte(plainResponse))
		}
	}

	path := randomString(12)
	server.HTTP.DELETE("/"+path, web.HTTPHandle(func(request web.Request, writer web.Writer) web.Response {
		return web.Response{}
	}), web.HandleOptions{})

	// HTML
	func() {
		req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path), nil)
		req.Header.Add("Accept", "text/html")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Network error: %s", err.Error())
		}
		if resp.StatusCode != 405 {
			t.Fatalf("Unexpected status code. Expected %d Got %d", 405, resp.StatusCode)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != htmlResponse {
			t.Fatalf("Unexpected body %v", body)
		}
	}()

	// JSON
	func() {
		req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path), nil)
		req.Header.Add("Accept", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Network error: %s", err.Error())
		}
		if resp.StatusCode != 405 {
			t.Fatalf("Unexpected status code. Expected %d Got %d", 405, resp.StatusCode)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != jsonResponse {
			t.Fatalf("Unexpected body %v", body)
		}
	}()

	// Plain
	func() {
		req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Network error: %s", err.Error())
		}
		if resp.StatusCode != 405 {
			t.Fatalf("Unexpected status code. Expected %d Got %d", 405, resp.StatusCode)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != plainResponse {
			t.Fatalf("Unexpected body %v", body)
		}
	}()
}

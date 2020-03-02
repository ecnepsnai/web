package web

import (
	"io/ioutil"
	"net/http"
	"strings"
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

func TestNotFoundHandle(t *testing.T) {
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
		req, _ := http.NewRequest("GET", "http://localhost:9557/"+randomString(12), nil)
		req.Header.Add("Accept", "text/html")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("Network error: %s", err.Error())
		}
		if resp.StatusCode != 404 {
			t.Errorf("Unexpected status code %d", resp.StatusCode)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != htmlResponse {
			t.Errorf("Unexpected body %v", body)
		}
	}()

	// JSON
	func() {
		req, _ := http.NewRequest("GET", "http://localhost:9557/"+randomString(12), nil)
		req.Header.Add("Accept", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("Network error: %s", err.Error())
		}
		if resp.StatusCode != 404 {
			t.Errorf("Unexpected status code %d", resp.StatusCode)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != jsonResponse {
			t.Errorf("Unexpected body %v", body)
		}
	}()

	// Plain
	func() {
		req, _ := http.NewRequest("GET", "http://localhost:9557/"+randomString(12), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("Network error: %s", err.Error())
		}
		if resp.StatusCode != 404 {
			t.Errorf("Unexpected status code %d", resp.StatusCode)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != plainResponse {
			t.Errorf("Unexpected body %v", body)
		}
	}()
}

func TestMethodNotAllowed(t *testing.T) {
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

	path := "/" + randomString(12)
	server.HTTP.DELETE(path, HTTPHandle(func(request Request, writer Writer) Response {
		return Response{}
	}), HandleOptions{})

	// HTML
	func() {
		req, _ := http.NewRequest("GET", "http://localhost:9557"+path, nil)
		req.Header.Add("Accept", "text/html")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("Network error: %s", err.Error())
		}
		if resp.StatusCode != 405 {
			t.Errorf("Unexpected status code %d", resp.StatusCode)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != htmlResponse {
			t.Errorf("Unexpected body %v", body)
		}
	}()

	// JSON
	func() {
		req, _ := http.NewRequest("GET", "http://localhost:9557"+path, nil)
		req.Header.Add("Accept", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("Network error: %s", err.Error())
		}
		if resp.StatusCode != 405 {
			t.Errorf("Unexpected status code %d", resp.StatusCode)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != jsonResponse {
			t.Errorf("Unexpected body %v", body)
		}
	}()

	// Plain
	func() {
		req, _ := http.NewRequest("GET", "http://localhost:9557"+path, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("Network error: %s", err.Error())
		}
		if resp.StatusCode != 405 {
			t.Errorf("Unexpected status code %d", resp.StatusCode)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != plainResponse {
			t.Errorf("Unexpected body %v", body)
		}
	}()
}

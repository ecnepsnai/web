package web_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"github.com/ecnepsnai/web"
)

func TestHTTPAddRoutes(t *testing.T) {
	t.Parallel()
	server := newServer()

	handle := func(w http.ResponseWriter, r web.Request) {
		w.WriteHeader(200)
	}
	options := web.HandleOptions{}

	server.HTTP.GET("/"+randomString(5), handle, options)
	server.HTTP.HEAD("/"+randomString(5), handle, options)
	server.HTTP.OPTIONS("/"+randomString(5), handle, options)
	server.HTTP.POST("/"+randomString(5), handle, options)
	server.HTTP.PUT("/"+randomString(5), handle, options)
	server.HTTP.PATCH("/"+randomString(5), handle, options)
	server.HTTP.DELETE("/"+randomString(5), handle, options)
}

func TestHTTPAuthenticated(t *testing.T) {
	t.Parallel()
	server := newServer()

	handle := func(w http.ResponseWriter, r web.Request) {
		w.WriteHeader(200)
	}
	authenticate := func(request *http.Request) interface{} {
		return 1
	}
	options := web.HandleOptions{
		AuthenticateMethod: authenticate,
	}

	path := randomString(5)

	server.HTTP.GET("/"+path, handle, options)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path))
	if err != nil {
		t.Fatalf("Network error: %s", err.Error())
	}
	if resp == nil {
		t.Fatalf("Nil response returned")
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Unexpected HTTP status code. Expected %d got %d", 200, resp.StatusCode)
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %s", err.Error())
	}
}

func TestHTTPUnauthenticated(t *testing.T) {
	t.Parallel()
	server := newServer()

	handle := func(w http.ResponseWriter, r web.Request) {
		w.WriteHeader(200)
	}
	authenticate := func(request *http.Request) interface{} {
		return nil
	}
	options := web.HandleOptions{
		AuthenticateMethod: authenticate,
	}

	path := randomString(5)

	server.HTTP.GET("/"+path, handle, options)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path))
	if err != nil {
		t.Fatalf("Network error: %s", err.Error())
	}
	if resp == nil {
		t.Fatalf("Nil response returned")
	}
	if resp.StatusCode != 401 {
		t.Fatalf("Unexpected HTTP status code. Expected %d got %d", 401, resp.StatusCode)
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %s", err.Error())
	}
}

func TestHTTPNotFound(t *testing.T) {
	t.Parallel()
	server := newServer()

	path := randomString(5)
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path))
	if err != nil {
		t.Fatalf("Network error: %s", err.Error())
	}
	if resp == nil {
		t.Fatalf("Nil response returned")
	}
	if resp.StatusCode != 404 {
		t.Fatalf("Unexpected HTTP status code. Expected %d got %d", 404, resp.StatusCode)
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %s", err.Error())
	}
}

func TestHTTPMethodNotAllowed(t *testing.T) {
	t.Parallel()
	server := newServer()

	handle := func(w http.ResponseWriter, r web.Request) {
		w.WriteHeader(200)
	}
	authenticate := func(request *http.Request) interface{} {
		return nil
	}
	options := web.HandleOptions{
		AuthenticateMethod: authenticate,
	}

	path := randomString(5)

	server.HTTP.POST("/"+path, handle, options)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path))
	if err != nil {
		t.Fatalf("Network error: %s", err.Error())
	}
	if resp == nil {
		t.Fatalf("Nil response returned")
	}
	if resp.StatusCode != 405 {
		t.Fatalf("Unexpected HTTP status code. Expected %d got %d", 405, resp.StatusCode)
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %s", err.Error())
	}
}

func TestHTTPHandleError(t *testing.T) {
	t.Parallel()
	server := newServer()

	handle := func(w http.ResponseWriter, r web.Request) {
		w.WriteHeader(403)
	}
	authenticate := func(request *http.Request) interface{} {
		return 1
	}
	options := web.HandleOptions{
		AuthenticateMethod: authenticate,
	}

	path := randomString(5)

	server.HTTP.GET("/"+path, handle, options)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path))
	if err != nil {
		t.Fatalf("Network error: %s", err.Error())
	}
	if resp == nil {
		t.Fatalf("Nil response returned")
	}
	if resp.StatusCode != 403 {
		t.Fatalf("Unexpected HTTP status code. Expected %d got %d", 403, resp.StatusCode)
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %s", err.Error())
	}
}

func TestHTTPResponse(t *testing.T) {
	t.Parallel()
	server := newServer()

	tmp := t.TempDir()
	data := randomString(5)
	name := randomString(5) + ".html"

	if err := os.WriteFile(path.Join(tmp, name), []byte(data), 0644); err != nil {
		t.Fatalf("Error making temporary file: %s", err.Error())
	}

	handle := func(w http.ResponseWriter, r web.Request) {
		f, err := os.Open(path.Join(tmp, name))
		if err != nil {
			t.Fatalf("Error opening temporary file: %s", err.Error())
		}
		defer f.Close()
		io.Copy(w, f)
	}
	options := web.HandleOptions{}

	path := randomString(5)

	server.HTTP.GET("/"+path, handle, options)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path))
	if err != nil {
		t.Fatalf("Network error: %s", err.Error())
	}
	if resp == nil {
		t.Fatalf("Nil response returned")
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Unexpected HTTP status code. Expected %d got %d", 200, resp.StatusCode)
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %s", err.Error())
	}
}

func TestHTTPContentType(t *testing.T) {
	t.Parallel()
	server := newServer()

	contentType := "application/amazing"
	handle := func(w http.ResponseWriter, r web.Request) {
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(200)
	}
	options := web.HandleOptions{}

	path := randomString(5)

	server.HTTP.GET("/"+path, handle, options)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path))
	if err != nil {
		t.Fatalf("Network error: %s", err.Error())
	}
	if resp == nil {
		t.Fatalf("Nil response returned")
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Unexpected HTTP status code. Expected %d got %d", 200, resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentType {
		t.Fatalf("Unexpected content type. Expected '%s' got '%s'", contentType, resp.Header.Get("Content-Type"))
	}
}

func TestHTTPHeaders(t *testing.T) {
	t.Parallel()
	server := newServer()

	headerKey := randomString(5)
	headerValue := randomString(5)
	handle := func(w http.ResponseWriter, r web.Request) {
		w.Header().Set(headerKey, headerValue)
		w.WriteHeader(200)
	}
	options := web.HandleOptions{}

	path := randomString(5)

	server.HTTP.GET("/"+path, handle, options)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path))
	if err != nil {
		t.Fatalf("Network error: %s", err.Error())
	}
	if resp == nil {
		t.Fatalf("Nil response returned")
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Unexpected HTTP status code. Expected %d got %d", 200, resp.StatusCode)
	}
	if resp.Header.Get(headerKey) != headerValue {
		t.Fatalf("Unexpected content type. Expected '%s: %s' got '%s: %s'", headerKey, headerValue, headerKey, resp.Header.Get(headerKey))
	}
}

func TestHTTPUnauthorizedMethod(t *testing.T) {
	t.Parallel()
	server := newServer()

	handle := func(w http.ResponseWriter, r web.Request) {
		w.WriteHeader(200)
	}
	authenticate := func(request *http.Request) interface{} {
		return nil
	}

	location := "somewhere-else"

	unauthorized := func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Location", location)
		w.WriteHeader(410)
	}
	options := web.HandleOptions{
		AuthenticateMethod: authenticate,
		UnauthorizedMethod: unauthorized,
	}

	path := randomString(5)

	server.HTTP.GET("/"+path, handle, options)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path))
	if err != nil {
		t.Fatalf("Network error: %s", err.Error())
	}
	if resp == nil {
		t.Fatalf("Nil response returned")
	}
	if resp.StatusCode != 410 {
		t.Fatalf("Unexpected HTTP status code. Expected %d got %d", 410, resp.StatusCode)
	}
	if resp.Header.Get("Location") != location {
		t.Fatalf("Missing expected HTTP header. Expected '%s' got '%s'", location, resp.Header.Get("Location"))
	}
}

func TestHTTPLargeBody(t *testing.T) {
	t.Parallel()
	server := newServer()

	handle := func(w http.ResponseWriter, r web.Request) {
		w.WriteHeader(200)
	}
	options := web.HandleOptions{
		MaxBodyLength: 10,
	}

	path := randomString(5)
	body := bytes.NewReader([]byte(randomString(50)))

	server.HTTP.POST("/"+path, handle, options)

	resp, err := http.Post(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path), "text-plain", body)
	if err != nil {
		t.Fatalf("Network error: %s", err.Error())
	}
	if resp == nil {
		t.Fatalf("Nil response returned")
	}
	if resp.StatusCode != 413 {
		t.Fatalf("Unexpected HTTP status code. Expected %d got %d", 413, resp.StatusCode)
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %s", err.Error())
	}
}

func TestHTTPRateLimit(t *testing.T) {
	t.Parallel()
	server := newServer()

	handle := func(w http.ResponseWriter, r web.Request) {
		w.WriteHeader(200)
	}
	options := web.HandleOptions{}

	path := randomString(5)

	server.Options.MaxRequestsPerSecond = 2
	server.HTTP.GET("/"+path, handle, options)

	testIdx := 1
	doTest := func(expectedStatus int) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path))
		if err != nil {
			t.Fatalf("Network error: %s", err.Error())
		}
		if resp.StatusCode != expectedStatus {
			t.Fatalf("Unexpected HTTP status code. Expected %d got %d in test %d", expectedStatus, resp.StatusCode, testIdx)
		}
		resp.Body.Close()
		testIdx++
	}

	doTest(200)
	time.Sleep(500 * time.Millisecond)
	doTest(200)
	time.Sleep(500 * time.Millisecond)
	doTest(200)
	time.Sleep(500 * time.Millisecond)
	doTest(200)
	doTest(200)
	doTest(429)
	time.Sleep(1 * time.Second)

	doTest(200)
	doTest(200)
	doTest(429)
}

func TestHTTPPreHandle(t *testing.T) {
	t.Parallel()
	server := newServer()

	path200 := randomString(5)
	path400 := randomString(5)

	handle := func(w http.ResponseWriter, r web.Request) {
		w.WriteHeader(200)
	}
	options := web.HandleOptions{
		PreHandle: func(w http.ResponseWriter, request *http.Request) error {
			if request.URL.Path == "/"+path400 {
				w.WriteHeader(400)
				return fmt.Errorf("boo")
			}
			return nil
		},
	}

	server.HTTP.GET("/"+path200, handle, options)
	server.HTTP.GET("/"+path400, handle, options)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path400))
	if err != nil {
		t.Fatalf("Network error: %s", err.Error())
	}
	if resp == nil {
		t.Fatalf("Nil response returned")
	}
	if resp.StatusCode != 400 {
		t.Fatalf("Unexpected HTTP status code. Expected %d got %d", 400, resp.StatusCode)
	}

	resp, err = http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path200))
	if err != nil {
		t.Fatalf("Network error: %s", err.Error())
	}
	if resp == nil {
		t.Fatalf("Nil response returned")
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Unexpected HTTP status code. Expected %d got %d", 200, resp.StatusCode)
	}
}

func TestHTTPPanic(t *testing.T) {
	t.Parallel()
	server := newServer()

	path := randomString(5)

	handle := func(w http.ResponseWriter, r web.Request) {
		panic("oh no!")
	}
	options := web.HandleOptions{}

	server.HTTP.GET("/"+path, handle, options)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path))
	if err != nil {
		t.Fatalf("Network error: %s", err.Error())
	}
	if resp == nil {
		t.Fatalf("Nil response returned")
	}
	if resp.StatusCode != 500 {
		t.Fatalf("Unexpected HTTP status code. Expected %d got %d", 500, resp.StatusCode)
	}
}

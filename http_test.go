package web

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"
)

func TestHTTPAddRoutes(t *testing.T) {
	handle := func(request Request) Response {
		return Response{}
	}
	options := HandleOptions{}

	path := randomString(5)
	server.HTTP.GET("/"+path, handle, options)
	server.HTTP.HEAD("/"+path, handle, options)
	server.HTTP.OPTIONS("/"+path, handle, options)
	server.HTTP.POST("/"+path, handle, options)
	server.HTTP.PUT("/"+path, handle, options)
	server.HTTP.PATCH("/"+path, handle, options)
	server.HTTP.DELETE("/"+path, handle, options)
}

func TestHTTPAuthenticated(t *testing.T) {
	handle := func(request Request) Response {
		return Response{}
	}
	authenticate := func(request *http.Request) interface{} {
		return 1
	}
	options := HandleOptions{
		AuthenticateMethod: authenticate,
	}

	path := randomString(5)

	server.HTTP.GET("/"+path, handle, options)

	resp, err := http.Get("http://localhost:9557/" + path)
	if err != nil {
		t.Errorf("Network error: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Errorf("Unexpected HTTP status code. Expected %d got %d", 200, resp.StatusCode)
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %s", err.Error())
	}
}

func TestHTTPUnauthenticated(t *testing.T) {
	handle := func(request Request) Response {
		return Response{}
	}
	authenticate := func(request *http.Request) interface{} {
		return nil
	}
	options := HandleOptions{
		AuthenticateMethod: authenticate,
	}

	path := randomString(5)

	server.HTTP.GET("/"+path, handle, options)

	resp, err := http.Get("http://localhost:9557/" + path)
	if err != nil {
		t.Errorf("Network error: %s", err.Error())
	}
	if resp.StatusCode != 401 {
		t.Errorf("Unexpected HTTP status code. Expected %d got %d", 401, resp.StatusCode)
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %s", err.Error())
	}
}

func TestHTTPNotFound(t *testing.T) {
	path := randomString(5)
	resp, err := http.Get("http://localhost:9557/" + path)
	if err != nil {
		t.Errorf("Network error: %s", err.Error())
	}
	if resp.StatusCode != 404 {
		t.Errorf("Unexpected HTTP status code. Expected %d got %d", 404, resp.StatusCode)
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %s", err.Error())
	}
}

func TestHTTPMethodNotAllowed(t *testing.T) {
	handle := func(request Request) Response {
		return Response{}
	}
	authenticate := func(request *http.Request) interface{} {
		return nil
	}
	options := HandleOptions{
		AuthenticateMethod: authenticate,
	}

	path := randomString(5)

	server.HTTP.POST("/"+path, handle, options)

	resp, err := http.Get("http://localhost:9557/" + path)
	if err != nil {
		t.Errorf("Network error: %s", err.Error())
	}
	if resp.StatusCode != 405 {
		t.Errorf("Unexpected HTTP status code. Expected %d got %d", 405, resp.StatusCode)
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %s", err.Error())
	}
}

func TestHTTPHandleError(t *testing.T) {
	handle := func(request Request) Response {
		return Response{
			Status: 403,
		}
	}
	authenticate := func(request *http.Request) interface{} {
		return 1
	}
	options := HandleOptions{
		AuthenticateMethod: authenticate,
	}

	path := randomString(5)

	server.HTTP.GET("/"+path, handle, options)

	resp, err := http.Get("http://localhost:9557/" + path)
	if err != nil {
		t.Errorf("Network error: %s", err.Error())
	}
	if resp.StatusCode != 403 {
		t.Errorf("Unexpected HTTP status code. Expected %d got %d", 403, resp.StatusCode)
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %s", err.Error())
	}
}

func TestHTTPResponse(t *testing.T) {
	data := randomString(5)
	name := randomString(5) + ".html"

	if err := ioutil.WriteFile(path.Join(tmpDir, name), []byte(data), 0644); err != nil {
		t.Errorf("Error making tempory file: %s", err.Error())
	}

	handle := func(request Request) Response {
		f, err := os.Open(path.Join(tmpDir, name))
		if err != nil {
			t.Errorf("Error opening temporary file: %s", err.Error())
		}
		return Response{
			Reader: f,
		}
	}
	options := HandleOptions{}

	path := randomString(5)

	server.HTTP.GET("/"+path, handle, options)

	resp, err := http.Get("http://localhost:9557/" + path)
	if err != nil {
		t.Errorf("Network error: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Errorf("Unexpected HTTP status code. Expected %d got %d", 200, resp.StatusCode)
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %s", err.Error())
	}
}

func TestHTTPContentType(t *testing.T) {
	contentType := "application/amazing"
	handle := func(request Request) Response {
		return Response{
			ContentType: contentType,
		}
	}
	options := HandleOptions{}

	path := randomString(5)

	server.HTTP.GET("/"+path, handle, options)

	resp, err := http.Get("http://localhost:9557/" + path)
	if err != nil {
		t.Errorf("Network error: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Errorf("Unexpected HTTP status code. Expected %d got %d", 200, resp.StatusCode)
	}
	if resp.Header.Get("Content-Type") != contentType {
		t.Errorf("Unexpected content type. Expected '%s' got '%s'", contentType, resp.Header.Get("Content-Type"))
	}
}

func TestHTTPHeaders(t *testing.T) {
	headerKey := randomString(5)
	headerValue := randomString(5)
	handle := func(request Request) Response {
		return Response{
			Headers: map[string]string{
				headerKey: headerValue,
			},
		}
	}
	options := HandleOptions{}

	path := randomString(5)

	server.HTTP.GET("/"+path, handle, options)

	resp, err := http.Get("http://localhost:9557/" + path)
	if err != nil {
		t.Errorf("Network error: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Errorf("Unexpected HTTP status code. Expected %d got %d", 200, resp.StatusCode)
	}
	if resp.Header.Get(headerKey) != headerValue {
		t.Errorf("Unexpected content type. Expected '%s: %s' got '%s: %s'", headerKey, headerValue, headerKey, resp.Header.Get(headerKey))
	}
}

func TestServeFile(t *testing.T) {
	data := randomString(5)
	name := randomString(5) + ".html"

	if err := ioutil.WriteFile(path.Join(tmpDir, name), []byte(data), 0644); err != nil {
		t.Errorf("Error making tempory file: %s", err.Error())
	}

	server.HTTP.Static("/static/*filepath", tmpDir)

	resp, err := http.Get("http://localhost:9557/static/" + name)
	if err != nil {
		t.Errorf("Network error: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Errorf("Unexpected HTTP status code. Expected %d got %d", 200, resp.StatusCode)
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %s", err.Error())
	}
}

func TestHTTPUnauthorizedMethod(t *testing.T) {
	handle := func(request Request) Response {
		return Response{}
	}
	authenticate := func(request *http.Request) interface{} {
		return nil
	}

	location := "somewhere-else"

	unauthorized := func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Location", location)
		w.WriteHeader(410)
	}
	options := HandleOptions{
		AuthenticateMethod: authenticate,
		UnauthorizedMethod: unauthorized,
	}

	path := randomString(5)

	server.HTTP.GET("/"+path, handle, options)

	resp, err := http.Get("http://localhost:9557/" + path)
	if err != nil {
		t.Errorf("Network error: %s", err.Error())
	}
	if resp.StatusCode != 410 {
		t.Errorf("Unexpected HTTP status code. Expected %d got %d", 410, resp.StatusCode)
	}
	if resp.Header.Get("Location") != location {
		t.Errorf("Missing expected HTTP header. Expected '%s' got '%s'", location, resp.Header.Get("Location"))
	}
}

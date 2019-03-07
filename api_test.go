package web

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAPIAddRoutes(t *testing.T) {
	handle := func(request Request) (interface{}, *Error) {
		return true, nil
	}
	options := HandleOptions{}

	path := randomString(5)
	server.API.GET("/"+path, handle, options)
	server.API.HEAD("/"+path, handle, options)
	server.API.OPTIONS("/"+path, handle, options)
	server.API.POST("/"+path, handle, options)
	server.API.PUT("/"+path, handle, options)
	server.API.PATCH("/"+path, handle, options)
	server.API.DELETE("/"+path, handle, options)
}

func TestAPIAuthenticated(t *testing.T) {
	handle := func(request Request) (interface{}, *Error) {
		return true, nil
	}
	authenticate := func(request *http.Request) interface{} {
		return 1
	}
	options := HandleOptions{
		AuthenticateMethod: authenticate,
	}

	path := randomString(5)

	server.API.GET("/"+path, handle, options)

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

func TestAPIUnauthenticated(t *testing.T) {
	handle := func(request Request) (interface{}, *Error) {
		return true, nil
	}
	authenticate := func(request *http.Request) interface{} {
		var object *string
		return object
	}
	options := HandleOptions{
		AuthenticateMethod: authenticate,
	}

	path := randomString(5)

	server.API.GET("/"+path, handle, options)

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

func TestAPINotFound(t *testing.T) {
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

func TestAPIMethodNotAllowed(t *testing.T) {
	handle := func(request Request) (interface{}, *Error) {
		return true, nil
	}
	authenticate := func(request *http.Request) interface{} {
		return nil
	}
	options := HandleOptions{
		AuthenticateMethod: authenticate,
	}

	path := randomString(5)

	server.API.POST("/"+path, handle, options)

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

func TestAPIHandleError(t *testing.T) {
	handle := func(request Request) (interface{}, *Error) {
		return nil, ValidationError("error")
	}
	authenticate := func(request *http.Request) interface{} {
		return 1
	}
	options := HandleOptions{
		AuthenticateMethod: authenticate,
	}

	path := randomString(5)

	server.API.GET("/"+path, handle, options)

	resp, err := http.Get("http://localhost:9557/" + path)
	if err != nil {
		t.Errorf("Network error: %s", err.Error())
	}
	if resp.StatusCode != 400 {
		t.Errorf("Unexpected HTTP status code. Expected %d got %d", 400, resp.StatusCode)
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %s", err.Error())
	}
}

func TestAPIUnauthorizedMethod(t *testing.T) {
	handle := func(request Request) (interface{}, *Error) {
		return true, nil
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

	server.API.GET("/"+path, handle, options)

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

func TestAPILargeBody(t *testing.T) {
	handle := func(request Request) (interface{}, *Error) {
		return true, nil
	}
	options := HandleOptions{
		MaxBodyLength: 10,
	}

	path := randomString(5)
	body := bytes.NewReader([]byte(randomString(50)))

	server.API.POST("/"+path, handle, options)

	resp, err := http.Post("http://localhost:9557/"+path, "text-plain", body)
	if err != nil {
		t.Errorf("Network error: %s", err.Error())
	}
	if resp.StatusCode != 413 {
		t.Errorf("Unexpected HTTP status code. Expected %d got %d", 413, resp.StatusCode)
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %s", err.Error())
	}
}

func TestAPIInvalidJSON(t *testing.T) {
	handle := func(request Request) (interface{}, *Error) {
		type exampleType struct {
			Foo string
			Bar string
		}

		example := exampleType{}
		if err := request.Decode(&example); err != nil {
			return nil, CommonErrors.BadRequest
		}
		return true, nil
	}
	options := HandleOptions{}

	path := randomString(5)
	body := bytes.NewReader([]byte(randomString(50)))

	server.API.POST("/"+path, handle, options)

	resp, err := http.Post("http://localhost:9557/"+path, "application/json", body)
	if err != nil {
		t.Errorf("Network error: %s", err.Error())
	}
	if resp.StatusCode != 400 {
		t.Errorf("Unexpected HTTP status code. Expected %d got %d", 400, resp.StatusCode)
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %s", err.Error())
	}
}

package api

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAddRoutes(t *testing.T) {
	handle := func(request Request) (interface{}, *Error) {
		return true, nil
	}
	options := HandleOptions{}

	server.APIGET("/", handle, options)
	server.APIHEAD("/", handle, options)
	server.APIOPTIONS("/", handle, options)
	server.APIPOST("/", handle, options)
	server.APIPUT("/", handle, options)
	server.APIPATCH("/", handle, options)
	server.APIDELETE("/", handle, options)
}

func TestAuthenticated(t *testing.T) {
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

	server.APIGET("/"+path, handle, options)

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

func TestUnauthenticated(t *testing.T) {
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

	server.APIGET("/"+path, handle, options)

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

func TestNotFound(t *testing.T) {
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

func TestMethodNotAllowed(t *testing.T) {
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

	server.APIPOST("/"+path, handle, options)

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

func TestHandleError(t *testing.T) {
	handle := func(request Request) (interface{}, *Error) {
		return nil, CommonErrors.Forbidden
	}
	authenticate := func(request *http.Request) interface{} {
		return 1
	}
	options := HandleOptions{
		AuthenticateMethod: authenticate,
	}

	path := randomString(5)

	server.APIGET("/"+path, handle, options)

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

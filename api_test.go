package web_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"testing"
	"time"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/web"
)

func TestAPIAddRoutes(t *testing.T) {
	t.Parallel()
	server := newServer()

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		return true, nil, nil
	}
	options := web.HandleOptions{}

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
	t.Parallel()
	server := newServer()

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		return true, nil, nil
	}
	authenticate := func(request *http.Request) interface{} {
		return 1
	}
	options := web.HandleOptions{
		AuthenticateMethod: authenticate,
	}

	path := randomString(5)

	server.API.GET("/"+path, handle, options)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path))
	if err != nil {
		t.Fatalf("Network error getting: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Unexpected HTTP status code. Expected %d got %d", 200, resp.StatusCode)
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %s", err.Error())
	}
}

func TestAPIUnauthenticated(t *testing.T) {
	t.Parallel()
	server := newServer()

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		return true, nil, nil
	}
	authenticate := func(request *http.Request) interface{} {
		var object *string
		return object
	}
	options := web.HandleOptions{
		AuthenticateMethod: authenticate,
	}

	path := randomString(5)

	server.API.GET("/"+path, handle, options)

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

func TestAPINotFound(t *testing.T) {
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

func TestAPIMethodNotAllowed(t *testing.T) {
	t.Parallel()
	server := newServer()

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		return true, nil, nil
	}
	authenticate := func(request *http.Request) interface{} {
		return nil
	}
	options := web.HandleOptions{
		AuthenticateMethod: authenticate,
	}

	path := randomString(5)

	server.API.POST("/"+path, handle, options)

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

func TestAPIHandleError(t *testing.T) {
	t.Parallel()
	server := newServer()

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		return nil, nil, web.ValidationError("error")
	}
	authenticate := func(request *http.Request) interface{} {
		return 1
	}
	options := web.HandleOptions{
		AuthenticateMethod: authenticate,
	}

	path := randomString(5)

	server.API.GET("/"+path, handle, options)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path))
	if err != nil {
		t.Fatalf("Network error: %s", err.Error())
	}
	if resp == nil {
		t.Fatalf("Nil response returned")
	}
	if resp.StatusCode != 400 {
		t.Fatalf("Unexpected HTTP status code. Expected %d got %d", 400, resp.StatusCode)
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %s", err.Error())
	}
}

func TestAPIHandlePanic(t *testing.T) {
	t.Parallel()
	server := newServer()

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		panic("oh no!")
	}
	authenticate := func(request *http.Request) interface{} {
		return 1
	}
	options := web.HandleOptions{
		AuthenticateMethod: authenticate,
	}

	path := randomString(5)

	server.API.GET("/"+path, handle, options)

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
	data := web.JSONResponse{}
	json.NewDecoder(resp.Body).Decode(&data)
	if data.Error == nil {
		t.Fatalf("No error in response body when one expected")
	}
}

func TestAPIUnauthorizedMethod(t *testing.T) {
	t.Parallel()
	server := newServer()

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		return true, nil, nil
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

	server.API.GET("/"+path, handle, options)

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

func TestAPILargeBody(t *testing.T) {
	t.Parallel()
	server := newServer()

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		return true, nil, nil
	}
	options := web.HandleOptions{
		MaxBodyLength: 10,
	}

	path := randomString(5)
	body := bytes.NewReader([]byte(randomString(50)))

	server.API.POST("/"+path, handle, options)

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

func TestAPIValidJSON(t *testing.T) {
	t.Parallel()
	server := newServer()

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		type exampleType struct {
			Foo string
			Bar string
		}

		example := exampleType{}
		if err := request.DecodeJSON(&example); err != nil {
			return nil, nil, web.CommonErrors.BadRequest
		}
		return true, nil, nil
	}
	options := web.HandleOptions{
		AuthenticateMethod: func(request *http.Request) interface{} {
			return true
		},
	}

	path := randomString(5)
	body := bytes.NewReader([]byte("{\"Foo\": \"1\", \"Bar\": \"2\"}"))

	server.API.POST("/"+path, handle, options)

	resp, err := http.Post(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path), "application/json", body)
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

func TestAPIInvalidJSON(t *testing.T) {
	t.Parallel()
	server := newServer()

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		type exampleType struct {
			Foo string
			Bar string
		}

		example := exampleType{}
		if err := request.DecodeJSON(&example); err != nil {
			return nil, nil, web.CommonErrors.BadRequest
		}
		return true, nil, nil
	}
	options := web.HandleOptions{}

	path := randomString(5)
	body := bytes.NewReader([]byte(randomString(50)))

	server.API.POST("/"+path, handle, options)

	resp, err := http.Post(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path), "application/json", body)
	if err != nil {
		t.Fatalf("Network error: %s", err.Error())
	}
	if resp == nil {
		t.Fatalf("Nil response returned")
	}
	if resp.StatusCode != 400 {
		t.Fatalf("Unexpected HTTP status code. Expected %d got %d", 400, resp.StatusCode)
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %s", err.Error())
	}
}

func TestAPIRateLimit(t *testing.T) {
	t.Parallel()
	server := newServer()

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		return true, nil, nil
	}
	options := web.HandleOptions{}

	path := randomString(5)

	server.Options.MaxRequestsPerSecond = 2
	server.API.GET("/"+path, handle, options)

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

func TestAPIResponse(t *testing.T) {
	t.Parallel()
	server := newServer()

	headerName := randomString(6)
	headerValue := randomString(6)
	cookieName := randomString(6)
	cookieValue := randomString(6)

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		return true, &web.APIResponse{
			Headers: map[string]string{
				headerName: headerValue,
			},
			Cookies: []http.Cookie{{
				Name:  cookieName,
				Value: cookieValue,
			}},
		}, nil
	}
	options := web.HandleOptions{}
	path := randomString(5)
	server.API.GET("/"+path, handle, options)

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

	if resp.Header.Get(headerName) != headerValue {
		t.Fatalf("Unexpected HTTP header. Expected %s got %s", headerValue, resp.Header.Get(headerName))
	}

	cookies := resp.Cookies()
	if len(cookies) != 1 {
		t.Fatalf("Unexpected number of cookies returned. Expected 1 got %d", len(cookies))
	}

	if cookies[0].Name != cookieName {
		t.Fatalf("Incorrect cookie name")
	}
	if cookies[0].Value != cookieValue {
		t.Fatalf("Incorrect cookie value")
	}
}

func TestAPILogLevel(t *testing.T) {
	logtic.Log.Reset()
	logFilePath := path.Join(t.TempDir(), "web.log")
	logtic.Log.FilePath = logFilePath

	stdout := &bytes.Buffer{}
	logtic.Log.Stdout = stdout
	logtic.Log.Stderr = stdout

	logtic.Log.Level = logtic.LevelDebug
	logtic.Log.Open()
	defer logtic.Log.Close()

	server := newServer()

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		return true, nil, nil
	}
	options := web.HandleOptions{}

	path := randomString(5)

	server.API.GET("/"+path, handle, options)

	http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path))
	server.Options.RequestLogLevel = logtic.LevelInfo
	http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path))

	logtic.Log.Close()
	debugPattern := regexp.MustCompile(`[0-9\-:TZ]+ \[DEBUG\]\[HTTP\] API Request: elapsed='[^']+' method='GET' remote_addr='[^']+' url='[^']+'`)
	infoPattern := regexp.MustCompile(`[0-9\-:TZ]+ \[INFO\]\[HTTP\] API Request: elapsed='[^']+' method='GET' remote_addr='[^']+' url='[^']+'`)
	f, err := os.OpenFile(logFilePath, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	logFileData, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	if !debugPattern.Match(logFileData) {
		t.Errorf("Did not find expected log line for API request\n----\n%s\n----", logFileData)
	}
	if !infoPattern.Match(logFileData) {
		t.Errorf("Did not find expected log line for API request\n----\n%s\n----", logFileData)
	}

	if stdout.Len() == 0 {
		t.Errorf("Did not find expected log line in stdout for API request")
	}

	logtic.Log.Reset()
	for _, arg := range os.Args {
		if arg == "-test.v=true" {
			logtic.Log.Level = logtic.LevelDebug
			logtic.Log.Open()
		}
	}
}

func TestAPIHandleNoLog(t *testing.T) {
	logtic.Log.Reset()
	logFilePath := path.Join(t.TempDir(), "web.log")
	logtic.Log.FilePath = logFilePath

	stdout := &bytes.Buffer{}
	logtic.Log.Stdout = stdout
	logtic.Log.Stderr = stdout

	logtic.Log.Level = logtic.LevelDebug
	logtic.Log.Open()
	defer logtic.Log.Close()

	server := newServer()

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		return true, nil, nil
	}

	path1 := randomString(5)
	path2 := randomString(5)

	server.API.GET("/"+path1, handle, web.HandleOptions{})
	server.API.GET("/"+path2, handle, web.HandleOptions{DontLogRequests: true})

	http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path1))
	http.Get(fmt.Sprintf("http://localhost:%d/%s", server.ListenPort, path2))

	logtic.Log.Close()
	path1Pattern := regexp.MustCompile(`[0-9\-:TZ]+ \[DEBUG\]\[HTTP\] API Request: elapsed='[^']+' method='GET' remote_addr='[^']+' url='/` + path1 + `'`)
	path2Pattern := regexp.MustCompile(`[0-9\-:TZ]+ \[DEBUG\]\[HTTP\] API Request: elapsed='[^']+' method='GET' remote_addr='[^']+' url='/` + path2 + `'`)
	f, err := os.OpenFile(logFilePath, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	logFileData, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	if !path1Pattern.Match(logFileData) {
		t.Errorf("Did not find expected log line for API request\n----\n%s\n----", logFileData)
	}
	if path2Pattern.Match(logFileData) {
		t.Errorf("Did not find expected log line for API request\n----\n%s\n----", logFileData)
	}

	if stdout.Len() == 0 {
		t.Errorf("Did not find expected log line in stdout for API request")
	}

	logtic.Log.Reset()
	for _, arg := range os.Args {
		if arg == "-test.v=true" {
			logtic.Log.Level = logtic.LevelDebug
			logtic.Log.Open()
		}
	}
}

func TestAPIPreHandle(t *testing.T) {
	t.Parallel()
	server := newServer()

	path200 := randomString(5)
	path400 := randomString(5)

	handle := func(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
		return true, nil, nil
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

	server.API.GET("/"+path200, handle, options)
	server.API.GET("/"+path400, handle, options)

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

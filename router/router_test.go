package router_test

import (
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/ecnepsnai/web/router"
)

func getListenAddress() string {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	l.Close()
	return l.Addr().String()
}

func testURL(t *testing.T, method, url string, expectedStatusCode int) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Unexpected status code for URL '%s'. Expected %d got %d", url, expectedStatusCode, resp.StatusCode)
	}
}

func TestRouterBasicGetIndex(t *testing.T) {
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

	testURL(t, "GET", "http://"+listenAddress+"/", 200)
}

func TestRouterBasicGetPath(t *testing.T) {
	t.Parallel()

	listenAddress := getListenAddress()

	server := router.New()
	server.Handle("GET", "/path", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	testURL(t, "GET", "http://"+listenAddress+"/path", 200)
}

func TestRouterBasicGetPathIndex(t *testing.T) {
	t.Parallel()

	listenAddress := getListenAddress()

	server := router.New()
	server.Handle("GET", "/path/", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	testURL(t, "GET", "http://"+listenAddress+"/path/", 200)
}

func TestRouterNotFound(t *testing.T) {
	t.Parallel()

	listenAddress := getListenAddress()

	server := router.New()
	server.Handle("GET", "/", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	server.Handle("GET", "/baz", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	testURL(t, "GET", "http://"+listenAddress+"/baz", 200)
	testURL(t, "GET", "http://"+listenAddress+"/foo", 404)
}

func TestRouterMethodNotAllowed(t *testing.T) {
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

	testURL(t, "POST", "http://"+listenAddress+"/", 405)
}

func TestRouterMultiplePaths(t *testing.T) {
	t.Parallel()

	listenAddress := getListenAddress()

	server := router.New()
	server.Handle("GET", "/one/two/three/", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	server.Handle("POST", "/one/two/three/", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	server.Handle("POST", "/one/two/four/", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	testURL(t, "GET", "http://"+listenAddress+"/one/two/three/", 200)
	testURL(t, "POST", "http://"+listenAddress+"/one/two/three/", 200)
	testURL(t, "POST", "http://"+listenAddress+"/one/two/four/", 200)
}

func TestRouterParameterizedPath(t *testing.T) {
	t.Parallel()

	listenAddress := getListenAddress()

	server := router.New()
	server.Handle("GET", "/one/:two/three/:four/", func(rw http.ResponseWriter, request router.Request) {
		p2 := request.Parameters["two"]
		if p2 != "two" {
			t.Errorf("Incorrect parameter value, Expected 'two' got '%s'", p2)
		}
		p4 := request.Parameters["four"]
		if p4 != "four" {
			t.Errorf("Incorrect parameter value, Expected 'four' got '%s'", p4)
		}
	})
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	testURL(t, "GET", "http://"+listenAddress+"/one/two/three/four/", 200)
}

func TestRouterAddInvalidMethod(t *testing.T) {
	t.Parallel()

	defer func() {
		recover()
	}()

	server := router.New()
	server.Handle("APPLESAUCE", "should/panic", func(rw http.ResponseWriter, request router.Request) {
		//
	})

	t.Errorf("No panic seen when one expected for adding invalid method")
}

func TestRouterAddReservedKeyword(t *testing.T) {
	t.Parallel()

	defer func() {
		recover()
	}()

	server := router.New()
	server.Handle("GET", "/__router_index", func(rw http.ResponseWriter, request router.Request) {
		//
	})

	t.Errorf("No panic seen when one expected for adding handle with reserved keyword")
}

func TestRouterAddInvalidPathNoSlash(t *testing.T) {
	t.Parallel()

	defer func() {
		recover()
	}()

	server := router.New()
	server.Handle("GET", "should/panic", func(rw http.ResponseWriter, request router.Request) {
		//
	})

	t.Errorf("No panic seen when one expected for adding path without leading slash")
}

func TestRouterAddInvalidPathParameterClash(t *testing.T) {
	t.Parallel()

	defer func() {
		recover()
	}()

	server := router.New()
	server.Handle("GET", "/one/:id/", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	server.Handle("GET", "/one/two/", func(rw http.ResponseWriter, request router.Request) {
		//
	})

	t.Errorf("No panic seen when one expected for adding path that collides with parametrized path")
}

func TestRouterAddInvalidDuplicateHandle(t *testing.T) {
	t.Parallel()

	defer func() {
		recover()
	}()

	server := router.New()
	server.Handle("GET", "/one/", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	server.Handle("GET", "/one/", func(rw http.ResponseWriter, request router.Request) {
		//
	})

	t.Errorf("No panic seen when one expected for adding duplicate handle")
}

func TestRouterRemoveHandle(t *testing.T) {
	t.Parallel()

	listenAddress := getListenAddress()

	server := router.New()
	server.Handle("GET", "/", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	server.Handle("GET", "/one/", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	server.Handle("GET", "/one/:id/", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	server.Handle("POST", "/one/:id/", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	server.Handle("GET", "/wildcard/*foo", func(rw http.ResponseWriter, r router.Request) {
		//
	})
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	testURL(t, "GET", "http://"+listenAddress+"/", 200)
	testURL(t, "GET", "http://"+listenAddress+"/one/", 200)
	testURL(t, "GET", "http://"+listenAddress+"/one/two/", 200)
	testURL(t, "GET", "http://"+listenAddress+"/wildcard/a/b/c/d", 200)

	server.RemoveHandle("GET", "/")
	server.RemoveHandle("GET", "/one/")
	server.RemoveHandle("GET", "/one/:id/")
	server.RemoveHandle("GET", "/wildcard/*foo")
	server.RemoveHandle("GET", "/foobar")
	server.RemoveHandle("GET", "")
	server.RemoveHandle("GET", "oops")

	testURL(t, "GET", "http://"+listenAddress+"/", 404)
	testURL(t, "GET", "http://"+listenAddress+"/one/", 404)
	testURL(t, "GET", "http://"+listenAddress+"/one/two/", 405)
	testURL(t, "POST", "http://"+listenAddress+"/one/two/", 200)
	testURL(t, "GET", "http://"+listenAddress+"/wildcard/a/b/c/d", 404)

	server.RemoveHandle("*", "*")

	testURL(t, "POST", "http://"+listenAddress+"/one/two/", 404)
}

func TestRouterNotFoundHandle(t *testing.T) {
	t.Parallel()

	listenAddress := getListenAddress()

	server := router.New()
	server.Handle("GET", "/cats", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	server.SetNotFoundHandle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("No dogs"))
	})
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	resp, err := http.Get("http://" + listenAddress + "/dogs")
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("Unexpected status code %d expected 404", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if string(body) != "No dogs" {
		t.Errorf("Incorrect response body")
	}
}

func TestRouterMethodNotAllowedHandle(t *testing.T) {
	t.Parallel()

	listenAddress := getListenAddress()

	server := router.New()
	server.Handle("POST", "/cats", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	server.SetMethodNotAllowedHandle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("No get"))
	})
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	resp, err := http.Get("http://" + listenAddress + "/cats")
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 405 {
		t.Errorf("Unexpected status code %d expected 405", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if string(body) != "No get" {
		t.Errorf("Incorrect response body: '%s'", body)
	}
}

func TestRouterWildcard(t *testing.T) {
	t.Parallel()

	listenAddress := getListenAddress()

	expectedUrl := "a/b/c/d/e/f/g/"

	server := router.New()
	server.Handle("GET", "/proxy/*url", func(rw http.ResponseWriter, r router.Request) {
		url := r.Parameters["url"]
		rw.Write([]byte(url))
	})
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	resp, err := http.Get("http://" + listenAddress + "/proxy/" + expectedUrl)
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	url := string(data)
	if url != expectedUrl {
		t.Errorf("Unexpected parameter value for wildcard. Expected '%s' got '%s'", expectedUrl, url)
	}
}

func TestRouterWildcardWithParameter(t *testing.T) {
	t.Parallel()

	listenAddress := getListenAddress()

	paramValue := "dogs"
	urlValue := "aregood"

	server := router.New()
	server.Handle("GET", "/req/:param/*url", func(rw http.ResponseWriter, r router.Request) {
		param := r.Parameters["param"]
		url := r.Parameters["url"]
		rw.Write([]byte(param + url))
	})
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	resp, err := http.Get("http://" + listenAddress + "/req/" + paramValue + "/" + urlValue)
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	url := string(data)
	expectedResponse := paramValue + urlValue
	if url != expectedResponse {
		t.Errorf("Unexpected parameter value for wildcard. Expected '%s' got '%s'", expectedResponse, url)
	}
}

func TestRouterWildcardRoot(t *testing.T) {
	t.Parallel()

	listenAddress := getListenAddress()

	expectedUrl := "a/b/c/d/e/f/g/"

	server := router.New()
	server.Handle("GET", "/*url", func(rw http.ResponseWriter, r router.Request) {
		url := r.Parameters["url"]
		rw.Write([]byte(url))
	})
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	resp, err := http.Get("http://" + listenAddress + "/" + expectedUrl)
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	url := string(data)
	if url != expectedUrl {
		t.Errorf("Unexpected parameter value for wildcard. Expected '%s' got '%s'", expectedUrl, url)
	}

	testURL(t, "POST", "http://"+listenAddress+"/"+expectedUrl, 405)
}

func TestRouterWildcardMultipleHandle(t *testing.T) {
	t.Parallel()

	server := router.New()
	server.Handle("GET", "/proxy/*url", func(rw http.ResponseWriter, r router.Request) {})
	server.Handle("HEAD", "/proxy/*url", func(rw http.ResponseWriter, r router.Request) {})
	server.Handle("DELETE", "/proxy/*url", func(rw http.ResponseWriter, r router.Request) {})
	server.Handle("GET", "/users/user/:username", func(rw http.ResponseWriter, r router.Request) {})
}

func TestRouterWildcardStaticSegmentClashA(t *testing.T) {
	t.Parallel()

	defer func() {
		recover()
	}()

	server := router.New()
	server.Handle("GET", "/proxy/*url", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	server.Handle("GET", "/proxy/roxy", func(rw http.ResponseWriter, request router.Request) {
		//
	})

	t.Errorf("No panic seen when one expected for adding path that collides with parametrized path")
}

func TestRouterWildcardStaticSegmentClashB(t *testing.T) {
	t.Parallel()

	defer func() {
		recover()
	}()

	server := router.New()
	server.Handle("GET", "/proxy/roxy", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	server.Handle("GET", "/proxy/*url", func(rw http.ResponseWriter, request router.Request) {
		//
	})

	t.Errorf("No panic seen when one expected for adding path that collides with parametrized path")
}

func TestRouterWildcardParameterNameClash(t *testing.T) {
	t.Parallel()

	defer func() {
		recover()
	}()

	server := router.New()
	server.Handle("GET", "/proxy/*url", func(rw http.ResponseWriter, request router.Request) {
		//
	})
	server.Handle("GET", "/proxy/*other", func(rw http.ResponseWriter, request router.Request) {
		//
	})

	t.Errorf("No panic seen when one expected for adding path that collides with wildcard path with different parameter name")
}

func TestRouterHandlePanic(t *testing.T) {
	t.Parallel()

	listenAddress := getListenAddress()

	server := router.New()
	server.Handle("GET", "/", func(rw http.ResponseWriter, request router.Request) {
		panic("aww fiddlesticks!")
	})
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	testURL(t, "GET", "http://"+listenAddress+"/", 500)
}

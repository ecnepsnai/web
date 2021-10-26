package router_test

import (
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"github.com/ecnepsnai/web/router"
)

func testStaticRequest(t *testing.T, method, url string, expectedStatus int, expectedMime string) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != expectedStatus {
		t.Errorf("Unexpected status code for URL '%s'. Expected %d got %d", url, expectedStatus, resp.StatusCode)
	}
	if resp.ContentLength == 0 {
		t.Errorf("Empty content for URL '%s'", url)
	}
	mime := resp.Header.Get("Content-Type")
	if mime != expectedMime {
		t.Errorf("Unexpected content type for URL '%s'. Expected '%s' got '%s'", url, expectedMime, mime)
	}
}

func TestRouterStaticSimple(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(path.Join(dir, "index.html"), []byte("foo"), os.ModePerm)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/static/assets/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/assets/index.html", 200, "text/html")
}

func TestRouterStatic(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	imageDir := path.Join(dir, "image")
	jsDir := path.Join(dir, "js")
	cssDir := path.Join(dir, "css")
	os.Mkdir(imageDir, os.ModePerm)
	os.Mkdir(jsDir, os.ModePerm)
	os.Mkdir(cssDir, os.ModePerm)
	os.WriteFile(path.Join(dir, "index.html"), []byte("foo"), os.ModePerm)
	os.WriteFile(path.Join(imageDir, "bg.jpg"), []byte("foo"), os.ModePerm)
	os.WriteFile(path.Join(jsDir, "main.js"), []byte("foo"), os.ModePerm)
	os.WriteFile(path.Join(cssDir, "style.css"), []byte("foo"), os.ModePerm)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/static/assets/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/assets/index.html", 200, "text/html")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/assets/image/bg.jpg", 200, "image/jpeg")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/assets/js/main.js", 200, "text/javascript")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/assets/css/style.css", 200, "text/css")
	testStaticRequest(t, "DELETE", "http://"+listenAddress+"/static/assets/index.html", 405, "text/plain; charset=utf-8")

	server.Stop()
	time.Sleep(5 * time.Millisecond)

	server = router.New()
	server.ServeFiles(dir, "/static/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/index.html", 200, "text/html")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/image/bg.jpg", 200, "image/jpeg")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/js/main.js", 200, "text/javascript")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/css/style.css", 200, "text/css")

	server.Stop()
	time.Sleep(5 * time.Millisecond)

	server = router.New()
	server.ServeFiles(dir, "/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	testStaticRequest(t, "GET", "http://"+listenAddress+"/index.html", 200, "text/html")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/image/bg.jpg", 200, "image/jpeg")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/js/main.js", 200, "text/javascript")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/css/style.css", 200, "text/css")
}

func TestRouterStaticPathTransversal(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	os.WriteFile(path.Join(dir, "index.html"), []byte("foo"), os.ModePerm)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/static")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/../../../../../../index.html", 200, "text/html")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/../../../../../../etc/password", 404, "text/plain; charset=utf-8")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/%2e/%2e/%2e/%2e/%2e/%2e/etc/password", 404, "text/plain; charset=utf-8")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/../etc/password", 404, "text/plain; charset=utf-8")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/..\\etc/password", 404, "text/plain; charset=utf-8")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/..\\/etc/password", 404, "text/plain; charset=utf-8")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/%2e%2e%2fetc/password", 404, "text/plain; charset=utf-8")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/%252e%252e%252fetc/password", 404, "text/plain; charset=utf-8")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/%c0%ae%c0%ae%c0%afetc/password", 404, "text/plain; charset=utf-8")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/..././etc/password", 404, "text/plain; charset=utf-8")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/...\\.\\etc/password", 404, "text/plain; charset=utf-8")
	testStaticRequest(t, "GET", "http://"+listenAddress+"/static/~/.bashrc", 404, "text/plain; charset=utf-8")
}

func TestRouterStaticAddReservedKeyword(t *testing.T) {
	t.Parallel()

	defer func() {
		recover()
	}()

	server := router.New()
	server.ServeFiles(t.TempDir(), "/__router_parameter")

	t.Errorf("No panic seen when one expected for adding handle with reserved keyword")
}

func testStaticHeadRequest(t *testing.T, url string, expectedMime string, expectedModifiedDate time.Time) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("Unexpected status code for URL '%s'. Expected %d got %d", url, 204, resp.StatusCode)
	}
	if resp.Header.Get("Content-Length") == "0" {
		t.Errorf("Empty content for URL '%s'", url)
	}
	mime := resp.Header.Get("Content-Type")
	if mime != expectedMime {
		t.Errorf("Unexpected content type for URL '%s'. Expected '%s' got '%s'", url, expectedMime, mime)
	}

	lastModified, err := time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", resp.Header.Get("Last-Modified"))
	if err != nil {
		t.Errorf("Invalid Last-Modified date for URL '%s': %s", url, err.Error())
	}

	if lastModified.Format(time.RFC3339) != expectedModifiedDate.Format(time.RFC3339) {
		t.Errorf("Invalid last modified date for URL '%s': Expected %s Got %s", url, expectedModifiedDate.String(), lastModified.String())
	}
}

func TestRouterStaticHeadRequest(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	imageDir := path.Join(dir, "image")
	jsDir := path.Join(dir, "js")
	cssDir := path.Join(dir, "css")
	indexTime := time.Now().UTC().AddDate(0, -1, -1)
	imageTime := time.Now().UTC().AddDate(0, -2, -1)
	jsTime := time.Now().UTC().AddDate(0, -3, -1)
	cssTime := time.Now().UTC().AddDate(0, -4, -1)
	os.Mkdir(imageDir, os.ModePerm)
	os.Mkdir(jsDir, os.ModePerm)
	os.Mkdir(cssDir, os.ModePerm)
	os.WriteFile(path.Join(dir, "index.html"), []byte("foo"), os.ModePerm)
	os.WriteFile(path.Join(imageDir, "bg.jpg"), []byte("foo"), os.ModePerm)
	os.WriteFile(path.Join(jsDir, "main.js"), []byte("foo"), os.ModePerm)
	os.WriteFile(path.Join(cssDir, "style.css"), []byte("foo"), os.ModePerm)
	os.Chtimes(path.Join(dir, "index.html"), time.Now(), indexTime)
	os.Chtimes(path.Join(imageDir, "bg.jpg"), time.Now(), imageTime)
	os.Chtimes(path.Join(jsDir, "main.js"), time.Now(), jsTime)
	os.Chtimes(path.Join(cssDir, "style.css"), time.Now(), cssTime)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/static/assets/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	testStaticHeadRequest(t, "http://"+listenAddress+"/static/assets/index.html", "text/html", indexTime)
	testStaticHeadRequest(t, "http://"+listenAddress+"/static/assets/image/bg.jpg", "image/jpeg", imageTime)
	testStaticHeadRequest(t, "http://"+listenAddress+"/static/assets/js/main.js", "text/javascript", jsTime)
	testStaticHeadRequest(t, "http://"+listenAddress+"/static/assets/css/style.css", "text/css", cssTime)
}

func TestRouterStaticIfModifiedSinceRequest(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	os.WriteFile(path.Join(dir, "index.html"), []byte("foo"), os.ModePerm)
	indexTime := time.Now().UTC().AddDate(0, -1, -1)
	os.Chtimes(path.Join(dir, "index.html"), time.Now(), indexTime)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/static/assets/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	url := "http://" + listenAddress + "/static/assets/index.html"
	expectedMime := "text/html"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("If-Modified-Since", time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("Unexpected status code for URL '%s'. Expected %d got %d", url, 204, resp.StatusCode)
	}
	if resp.Header.Get("Content-Length") == "0" {
		t.Errorf("Empty content for URL '%s'", url)
	}
	mime := resp.Header.Get("Content-Type")
	if mime != expectedMime {
		t.Errorf("Unexpected content type for URL '%s'. Expected '%s' got '%s'", url, expectedMime, mime)
	}

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("If-Modified-Since", time.Now().UTC().AddDate(-1, 0, 0).Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Unexpected status code for URL '%s'. Expected %d got %d", url, 200, resp.StatusCode)
	}
	if resp.Header.Get("Content-Length") == "0" {
		t.Errorf("Empty content for URL '%s'", url)
	}
	mime = resp.Header.Get("Content-Type")
	if mime != expectedMime {
		t.Errorf("Unexpected content type for URL '%s'. Expected '%s' got '%s'", url, expectedMime, mime)
	}

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("If-Modified-Since", "foobar")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("Unexpected status code for URL '%s'. Expected %d got %d", url, 204, resp.StatusCode)
	}
	if resp.Header.Get("Content-Length") == "0" {
		t.Errorf("Empty content for URL '%s'", url)
	}
	mime = resp.Header.Get("Content-Type")
	if mime != expectedMime {
		t.Errorf("Unexpected content type for URL '%s'. Expected '%s' got '%s'", url, expectedMime, mime)
	}
}

func TestRouterStaticCacheControlHeader(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	os.WriteFile(path.Join(dir, "index.html"), []byte("foo"), os.ModePerm)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/static/assets/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)

	url := "http://" + listenAddress + "/static/assets/index.html"
	expectedMime := "text/html"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Unexpected status code for URL '%s'. Expected %d got %d", url, 200, resp.StatusCode)
	}
	mime := resp.Header.Get("Content-Type")
	if mime != expectedMime {
		t.Errorf("Unexpected content type for URL '%s'. Expected '%s' got '%s'", url, expectedMime, mime)
	}
	cacheControl := resp.Header.Get("Cache-Control")
	expectedCacheControl := "max-age=86400; public"
	if cacheControl != expectedCacheControl {
		t.Errorf("Unexpected cache control for URL '%s'. Expected '%s' got '%s'", url, expectedCacheControl, cacheControl)
	}

	router.CacheMaxAge = 0

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Unexpected status code for URL '%s'. Expected %d got %d", url, 200, resp.StatusCode)
	}
	mime = resp.Header.Get("Content-Type")
	if mime != expectedMime {
		t.Errorf("Unexpected content type for URL '%s'. Expected '%s' got '%s'", url, expectedMime, mime)
	}
	cacheControl = resp.Header.Get("Cache-Control")
	if cacheControl != "" {
		t.Errorf("Unexpected cache control for URL '%s'.", url)
	}
}

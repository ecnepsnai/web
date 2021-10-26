package router

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func defaultNotFoundHandle(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(404)

	accept := strings.ToLower(req.Header.Get("Accept"))
	if strings.Contains(accept, "html") {
		body := []byte("<html><body><h1>404 Not Found</h1></body></html>")
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
		w.Header().Set("Date", timeToHTTPDate(time.Now().UTC()))
		w.Write(body)
		return
	}

	body := []byte("404 not found")
	w.Header().Add("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
	w.Header().Set("Date", timeToHTTPDate(time.Now().UTC()))
	w.Write(body)
}

func defaultMethodNotAllowedHandle(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(405)

	accept := strings.ToLower(req.Header.Get("Accept"))
	if strings.Contains(accept, "html") {
		body := []byte("<html><body><h1>405 Method Not Allowed</h1></body></html>")
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
		w.Header().Set("Date", timeToHTTPDate(time.Now().UTC()))
		w.Write(body)
		return
	}

	body := []byte("405 method not allowed")
	w.Header().Add("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
	w.Header().Set("Date", timeToHTTPDate(time.Now().UTC()))
	w.Write(body)
}

package web_test

import (
	"net/http"
	"testing"

	"github.com/ecnepsnai/web"
)

func TestRealRemoteAddr(t *testing.T) {
	requestWithHeader := func(key, value string) *http.Request {
		r := &http.Request{
			Header: http.Header{},
		}
		r.Header.Set(key, value)
		return r
	}

	if ip := web.RealRemoteAddr(requestWithHeader("X-Real-IP", "127.0.0.1")).String(); ip != "127.0.0.1" {
		t.Errorf("Unexpected result from RealRemoteAddr: expected '%s' got '%s'", "127.0.0.1", ip)
	}
	if ip := web.RealRemoteAddr(requestWithHeader("X-Forwarded-For", "127.0.0.2")).String(); ip != "127.0.0.2" {
		t.Errorf("Unexpected result from RealRemoteAddr: expected '%s' got '%s'", "127.0.0.2", ip)
	}
	if ip := web.RealRemoteAddr(requestWithHeader("CF-Connecting-IP", "127.0.0.3")).String(); ip != "127.0.0.3" {
		t.Errorf("Unexpected result from RealRemoteAddr: expected '%s' got '%s'", "127.0.0.3", ip)
	}
	if ip := web.RealRemoteAddr(requestWithHeader("X-Real-IP", "1::1")).String(); ip != "1::1" {
		t.Errorf("Unexpected result from RealRemoteAddr: expected '%s' got '%s'", "1::1", ip)
	}
	if ip := web.RealRemoteAddr(requestWithHeader("X-Forwarded-For", "1::2")).String(); ip != "1::2" {
		t.Errorf("Unexpected result from RealRemoteAddr: expected '%s' got '%s'", "1::2", ip)
	}
	if ip := web.RealRemoteAddr(requestWithHeader("CF-Connecting-IP", "1::3")).String(); ip != "1::3" {
		t.Errorf("Unexpected result from RealRemoteAddr: expected '%s' got '%s'", "1::3", ip)
	}

	r := &http.Request{
		Header:     http.Header{},
		RemoteAddr: "127.0.0.4:1234",
	}
	if ip := web.RealRemoteAddr(r).String(); ip != "127.0.0.4" {
		t.Errorf("Unexpected result from RealRemoteAddr: expected '%s' got '%s'", "127.0.0.4", ip)
	}

	r = &http.Request{
		Header:     http.Header{},
		RemoteAddr: "[1::4]:1234",
	}
	if ip := web.RealRemoteAddr(r).String(); ip != "1::4" {
		t.Errorf("Unexpected result from RealRemoteAddr: expected '%s' got '%s'", "1::4", ip)
	}
}

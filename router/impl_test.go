package router

import (
	"net/http"
	"testing"
)

func TestImpl(t *testing.T) {
	s := New()
	s.Handle("GET", "/", func(rw http.ResponseWriter, r Request) { /* */ })
	s.Handle("GET", "/dir", func(rw http.ResponseWriter, r Request) { /* */ })
	s.Handle("GET", "/dir/", func(rw http.ResponseWriter, r Request) { /* */ })
	s.Handle("GET", "/dir/file", func(rw http.ResponseWriter, r Request) { /* */ })

	if len(s.impl.Index.Methods) > 0 {
		t.Errorf("No methods should be present on the index")
	}
	if len(s.impl.Index.Children) != 2 {
		t.Errorf("Incorrect number of children on the index")
	}

	// Asset that these exact children are available, this will easily panic if they are not
	s.impl.Index.Children[pathKeyIndex].Methods["GET"](nil, Request{})
	s.impl.Index.Children["dir"].Methods["GET"](nil, Request{})
	s.impl.Index.Children["dir"].Children[pathKeyIndex].Methods["GET"](nil, Request{})
	s.impl.Index.Children["dir"].Children["file"].Methods["GET"](nil, Request{})
}

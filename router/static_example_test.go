package router_test

import (
	"github.com/ecnepsnai/web/router"
)

func ExampleServer_ServeFiles() {
	server := router.New()

	localDirectory := "/usr/share/http" // Directory to serve files from
	urlRoot := "/assets/"               // Top-level path for all requests to be directed to the local directory

	server.ServeFiles(localDirectory, urlRoot)
	// Now any HTTP GET or HEAD requests to /assets/ will read files from /usr/share/http.
	// For example:
	// HTTP GET "/assets/index.html" will read file "/usr/share/http/index.html"
}

/*
Package web is a HTTP server for Golang applications.

It is suitable for both front-end and back-end use, being able to deliver static content, act as a REST-ful JSON server,
and as a WebSocket server.

It includes simple controls to allow for user authentication with contextual data being available in every request, and
provides simple per-user rate-limiting.
*/
package web

import "github.com/ecnepsnai/logtic"

var log = logtic.Log.Connect("HTTP")

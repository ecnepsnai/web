/*
Package web is a full-featured HTTP router and server for Go applications, suitable for serving static files, REST APIs,
and more.

Web includes these features:
  - HTTP range support
  - Static file serving
  - Directory listings
  - Websockets
  - Per-IP rate limiting
  - Per-request contextual data

Web offers four APIs for developers to choose from:

# API

API provides everything you need to build poweful REST APIs using JSON. Define your routes and easily accept and return
data as JSON.

Example:

	router := web.New("[::]:8080")
	router.API.Get("/users", getUsers, options)
	router.API.Get("/users/:username", getUsers, options)

For more information, see the documentation of [web.API].

# HTTPEasy

HTTPEasy provides a straightforward interface to accept HTTP requets and return data.

Example:

	router := web.New("[::]:8080")
	router.HTTPEasy.Get("/index.html", getIndex, options)
	router.HTTPEasy.Get("/cat.jpg", getKitty, options)

For more information, see the documentation of [web.HTTPEasy].

# HTTP

HTTP provides full access to the original HTTP request, allowing you total control over the response, whatever that may be.

Example:

	router := web.New("[::]:8080")
	router.HTTP.Get("/index.html", getIndex, options)
	router.HTTP.Get("/cat.jpg", getKitty, options)

For more information, see the documentation of [web.HTTP].

# Websockets

This package also provides a wrapper for [github.com/gorilla/websocket]

Example:

	router := web.New("[::]:8080")
	router.Socket("/ws", handleSocket, options)
*/
package web

import "github.com/ecnepsnai/logtic"

var log = logtic.Log.Connect("HTTP")

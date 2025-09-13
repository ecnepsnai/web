> [!IMPORTANT]  
> This package is only receiving bug and security fixes - No new features will be added. The evolution of this package is [git.ecn.io/ian/w3](https://git.ecn.io/ian/w3)

# Web

[![Go Report Card](https://goreportcard.com/badge/github.com/ecnepsnai/web?style=flat-square)](https://goreportcard.com/report/github.com/ecnepsnai/web)
[![Godoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/ecnepsnai/web)
[![Releases](https://img.shields.io/github/release/ecnepsnai/web/all.svg?style=flat-square)](https://github.com/ecnepsnai/web/releases)
[![LICENSE](https://img.shields.io/github/license/ecnepsnai/web.svg?style=flat-square)](https://github.com/ecnepsnai/web/blob/master/LICENSE)

The web project provides two packages, web and router.

## Web

Package web is a HTTP server for Golang applications.

It is suitable for both front-end and back-end use, being able to deliver static content, act as a REST-ful JSON server,
and as a WebSocket server.

It includes simple controls to allow for user authentication with contextual data being available in every request, and
provides simple per-user rate-limiting.

## Router

Package router provides a simple & efficient parametrized HTTP router.

A HTTP router allows you to map a HTTP request method and path to a specific function. A parameterized HTTP router
allows you to designate specific portions of the request path as a parameter, which can later be fetched during the
request itself.

This package allows you modify the routing table ad-hoc, even while the server is running.

# Documentation & Examples

For full documentation including examples please see the official [package documentation](https://pkg.go.dev/github.com/ecnepsnai/web)












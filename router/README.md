# router

Package router provides a simple & efficient parametrized HTTP router.

A HTTP router allows you to map a HTTP request method and path to a specific function. A parameterized HTTP router
allows you to designate specific portions of the request path as a parameter, which can later be fetched during the
request itself.

This package allows you modify the routing table ad-hoc, even while the server is running.

# Install

This package is provided by 'github.com/ecnepsnai/web', so add that to your go.mod file:

```
go get github.com/ecnepsnai/web
```
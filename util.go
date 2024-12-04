package web

import (
	"net"
	"net/http"
)

// RealRemoteAddr will try to get the real IP address of the incoming connection taking proxies into
// consideration. This function looks for the `X-Real-IP`, `X-Forwarded-For`, and `CF-Connecting-IP`
// headers, and if those don't exist will return the remote address of the connection.
//
// Will never return nil, if it is unable to get a valid address it will return 0.0.0.0
func RealRemoteAddr(r *http.Request) net.IP {
	if ip := net.ParseIP(r.Header.Get("X-Real-IP")); ip != nil {
		return ip
	}
	if ip := net.ParseIP(r.Header.Get("X-Forwarded-For")); ip != nil {
		return ip
	}
	if ip := net.ParseIP(r.Header.Get("CF-Connecting-IP")); ip != nil {
		return ip
	}

	ipStr, _, _ := net.SplitHostPort(r.RemoteAddr)
	if ip := net.ParseIP(ipStr); ip != nil {
		return ip
	}

	return net.IPv4(0, 0, 0, 0)
}

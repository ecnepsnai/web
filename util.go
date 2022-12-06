package web

import (
	"net"
	"net/http"
)

func getRealIP(r *http.Request) net.IP {
	if ip := net.ParseIP(r.Header.Get("X-Real-IP")); ip != nil {
		return ip
	}
	if ip := net.ParseIP(r.Header.Get("X-Forwarded-For")); ip != nil {
		return ip
	}
	if ip := net.ParseIP(getIPFromRemoteAddr(r.RemoteAddr)); ip != nil {
		return ip
	}

	return net.IPv4(0, 0, 0, 0)
}

// getIPFromRemoteAddr strip the port from a socket address (address:port, return address). Also unwraps IPv6 addresses.
func getIPFromRemoteAddr(remoteAddr string) string {
	addr := stripPortFromSocketAddr(remoteAddr)
	if addr[0] == '[' && addr[len(addr)-1] == ']' {
		return addr[1 : len(addr)-1]
	}

	return addr
}

// stripPortFromSocketAddr strip the port from a socket address (address:port, return address)
func stripPortFromSocketAddr(socketAddr string) string {
	length := len(socketAddr)
	for i := length - 1; i >= 0; i-- {
		if socketAddr[i] == ':' {
			return socketAddr[0:i]
		}
	}

	return socketAddr
}

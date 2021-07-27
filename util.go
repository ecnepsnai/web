package web

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
	if socketAddr[length-6] == ':' {
		return socketAddr[0 : length-6]
	} else if socketAddr[length-5] == ':' {
		return socketAddr[0 : length-5]
	} else if socketAddr[length-4] == ':' {
		return socketAddr[0 : length-4]
	} else if socketAddr[length-3] == ':' {
		return socketAddr[0 : length-3]
	} else if socketAddr[length-2] == ':' {
		return socketAddr[0 : length-2]
	}

	return socketAddr
}

package web

// getIPFromRemoteAddr strip the port from a socket address (address:port, return address)
func getIPFromRemoteAddr(remoteAddr string) string {
	length := len(remoteAddr)
	if remoteAddr[length-6] == ':' {
		return remoteAddr[0 : length-6]
	} else if remoteAddr[length-5] == ':' {
		return remoteAddr[0 : length-5]
	} else if remoteAddr[length-4] == ':' {
		return remoteAddr[0 : length-4]
	} else if remoteAddr[length-3] == ':' {
		return remoteAddr[0 : length-3]
	} else if remoteAddr[length-2] == ':' {
		return remoteAddr[0 : length-2]
	}

	return remoteAddr
}

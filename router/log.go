package router

import (
	"bytes"

	"github.com/ecnepsnai/logtic"
)

// mute the following unhelpful events from log lines

var logMutePatterns = [][]byte{
	[]byte("http: TLS handshake error from"),
}

type muteLogger struct {
	source *logtic.Source
	level  int
}

func (l muteLogger) Write(p []byte) (n int, err error) {
	length := len(p)

	mute := false
	for _, pattern := range logMutePatterns {
		if bytes.Contains(p, pattern) {
			mute = true
			break
		}
	}
	if !mute {
		l.source.Write(l.level, string(p[0:length-1]))
	}

	return len(p), nil
}

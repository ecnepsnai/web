package api

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"testing"
)

var Verbose bool
var server *Server

func isVerbose() bool {
	for _, arg := range os.Args {
		if arg == "-test.v=true" {
			return true
		}
	}

	return false
}

func testSetup() {
	server = New("127.0.0.1:9557")
	go func() {
		if err := server.Start(); err != nil {
			panic(err)
		}
	}()
}

func testTeardown() {
	server.Stop()
}

func TestMain(m *testing.M) {
	Verbose = isVerbose()
	testSetup()
	retCode := m.Run()
	testTeardown()
	os.Exit(retCode)
}

func randomString(length uint16) string {
	randB := make([]byte, length)
	rand.Read(randB)
	return hex.EncodeToString(randB)
}

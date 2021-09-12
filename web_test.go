package web_test

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/web"
)

var serverLock = &sync.Mutex{}
var servers = []*web.Server{}

func newServer() *web.Server {
	server := web.New(":0")
	serverLock.Lock()
	servers = append(servers, server)
	serverLock.Unlock()
	go server.Start()

	// It can take a couple cycles for the server to be ready, so wait for the port to be populated before returning
	i := 0
	for i < 10 {
		if server.ListenPort > 0 {
			break
		}
		i++
		time.Sleep(5 * time.Millisecond)
	}
	if server.ListenPort == 0 {
		panic("Server didn't start in time")
	}

	return server
}

func testSetup() {
	for _, arg := range os.Args {
		if arg == "-test.v=true" {
			logtic.Log.Level = logtic.LevelDebug
			logtic.Log.Open()
		}
	}
}

func testTeardown() {
	for _, server := range servers {
		go server.Stop()
	}
}

func TestMain(m *testing.M) {
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

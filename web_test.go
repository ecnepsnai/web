package web

import (
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/ecnepsnai/logtic"
)

var Verbose bool
var server *Server
var tmpDir string
var logFile *logtic.File

func isTestVerbose() bool {
	for _, arg := range os.Args {
		if arg == "-test.v=true" {
			return true
		}
	}

	return false
}

func testSetup() {
	t, err := ioutil.TempDir("", "web")
	if err != nil {
		panic(err)
	}
	tmpDir = t

	if Verbose {
		initLogtic()
	}

	server = New("127.0.0.1:9557")
	testStartServer()
}

func testStartServer() {
	go func() {
		if err := server.Start(); err != nil {
			panic(err)
		}
	}()
}

func testTeardown() {
	server.Stop()
	os.RemoveAll(tmpDir)
	if logFile != nil {
		logFile.Close()
	}
}

func initLogtic() {
	file, _, err := logtic.New(path.Join(tmpDir, "web.log"), logtic.LevelDebug, "test runner")
	if err != nil {
		panic(err)
	}
	logFile = file
}

func TestMain(m *testing.M) {
	Verbose = isTestVerbose()
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

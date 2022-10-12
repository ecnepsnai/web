package router_test

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"math/big"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/web/router"
)

func TestMuteLogger(t *testing.T) {
	b := &bytes.Buffer{}
	logtic.Log.Open()
	logtic.Log.Level = logtic.LevelDebug
	logtic.Log.Stdout = b
	logtic.Log.Stderr = b

	listenAddress := getListenAddress()

	var pKey crypto.PrivateKey
	var err error
	pKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		panic(err)
	}

	pub := pKey.(crypto.Signer).Public()
	tpl := &x509.Certificate{
		SerialNumber:          &big.Int{},
		NotBefore:             time.Now().UTC().AddDate(-100, 0, 0),
		NotAfter:              time.Now().UTC().AddDate(100, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment,
		BasicConstraintsValid: true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, tpl, tpl, pub, pKey)
	if err != nil {
		panic(err)
	}

	l, err := tls.Listen("tcp", listenAddress, &tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{certBytes},
				PrivateKey:  pKey,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	server := router.New()
	server.ServeFiles(t.TempDir(), "/")
	go func() {
		server.Serve(l)
	}()
	time.Sleep(5 * time.Millisecond)

	c, err := net.Dial("tcp", listenAddress)
	if err != nil {
		panic(err)
	}
	if _, err := c.Write([]byte("")); err != nil {
		panic(err)
	}
	c.Close()
	time.Sleep(5 * time.Millisecond)

	output := b.String()
	if strings.Contains(output, "http: TLS handshake error from") {
		t.Errorf("log output contains forbidden log events")
	}
}

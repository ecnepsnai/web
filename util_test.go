package web

import "testing"

func TestGetIPFromRemoteAddr(t *testing.T) {
	if ip := getIPFromRemoteAddr("127.0.0.1:42336"); ip != "127.0.0.1" {
		t.Errorf("Incorrect result for IP address. Expected: %s Actual: %s", "127.0.0.1", ip)
	}

	if ip := getIPFromRemoteAddr("127.0.0.1:4233"); ip != "127.0.0.1" {
		t.Errorf("Incorrect result for IP address. Expected: %s Actual: %s", "127.0.0.1", ip)
	}

	if ip := getIPFromRemoteAddr("[1::1]:4233"); ip != "1::1" {
		t.Errorf("Incorrect result for IP address. Expected: %s Actual: %s", "1::1", ip)
	}
}

func BenchmarkGetIPFromRemoteAddr(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getIPFromRemoteAddr("[1::1]:42336")
	}
}

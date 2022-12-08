package router_test

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"github.com/ecnepsnai/web/router"
)

var sampleData []byte

func init() {
	sampleData = make([]byte, 500)

	i := 0
	y := 1
	for y <= 5 {
		d := []byte(fmt.Sprintf("%d", y))[0]
		for x := 0; x < 100; x++ {
			sampleData[i] = d
			i++
		}
		y++
	}
}

func TestRangeHEADRequest(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(path.Join(dir, "data.txt"), sampleData, os.ModePerm)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)
	url := "http://" + listenAddress + "/data.txt"

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", "rangetest/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("unexpected HTTP status code. Expected %d got %d", 200, resp.StatusCode)
	}

	if resp.ContentLength != int64(len(sampleData)) {
		t.Fatalf("incorrect value of content length header. Expected %d got %d", len(sampleData), resp.ContentLength)
	}

	ct, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		panic(err)
	}
	if ct != "text/plain" {
		t.Fatalf("incorrect value of content typ header. Expected 'text/plain' got '%s'", ct)
	}

	if value := resp.Header.Get("Accept-Ranges"); value != "bytes" {
		t.Fatalf("incorrect or missing Accept-Ranges header. Expected bytes got %s", value)
	}
}

func TestRangeGetAllDataWithoutRange(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(path.Join(dir, "data.txt"), sampleData, os.ModePerm)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)
	url := "http://" + listenAddress + "/data.txt"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", "rangetest/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("unexpected HTTP status code. Expected %d got %d", 200, resp.StatusCode)
	}

	if resp.ContentLength != int64(len(sampleData)) {
		t.Fatalf("incorrect value of content length header. Expected %d got %d", len(sampleData), resp.ContentLength)
	}

	ct, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		panic(err)
	}
	if ct != "text/plain" {
		t.Fatalf("incorrect value of content typ header. Expected 'text/plain' got '%s'", ct)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if len(respData) != len(sampleData) {
		t.Fatalf("incorrect data length. Expected %d got %d", len(sampleData), len(respData))
	}

	if !bytes.Equal(respData, sampleData) {
		t.Fatalf("invalid data returned")
	}
}

func TestRangeGetAllDataWithRange(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(path.Join(dir, "data.txt"), sampleData, os.ModePerm)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)
	url := "http://" + listenAddress + "/data.txt"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", "rangetest/1.0")
	req.Header.Add("Range", "bytes=0-")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 206 {
		t.Fatalf("unexpected HTTP status code. Expected %d got %d", 206, resp.StatusCode)
	}

	if resp.ContentLength != int64(len(sampleData)) {
		t.Fatalf("incorrect value of content length header. Expected %d got %d", len(sampleData), resp.ContentLength)
	}

	ct, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		panic(err)
	}
	if ct != "text/plain" {
		t.Fatalf("incorrect value of content typ header. Expected 'text/plain' got '%s'", ct)
	}

	if value := resp.Header.Get("Content-Range"); value != "bytes 0-499/500" {
		t.Fatalf("incorrect or missing Accept-Ranges header. Expected 'bytes 0-499/500' got '%s'", value)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if len(respData) != len(sampleData) {
		t.Fatalf("incorrect data length. Expected %d got %d", len(sampleData), len(respData))
	}

	if !bytes.Equal(respData, sampleData) {
		t.Fatalf("invalid data returned")
	}
}

func TestRangeGetSingleAbsoluteRange(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(path.Join(dir, "data.txt"), sampleData, os.ModePerm)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)
	url := "http://" + listenAddress + "/data.txt"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", "rangetest/1.0")
	req.Header.Add("Range", "bytes=0-99")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 206 {
		t.Fatalf("unexpected HTTP status code. Expected %d got %d", 206, resp.StatusCode)
	}

	if resp.ContentLength != int64(100) {
		t.Fatalf("incorrect value of content length header. Expected %d got %d", 100, resp.ContentLength)
	}

	ct, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		panic(err)
	}
	if ct != "text/plain" {
		t.Fatalf("incorrect value of content typ header. Expected 'text/plain' got '%s'", ct)
	}

	if value := resp.Header.Get("Content-Range"); value != "bytes 0-99/500" {
		t.Fatalf("incorrect or missing Accept-Ranges header. Expected 'bytes 0-99/500' got '%s'", value)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if len(respData) != 100 {
		t.Fatalf("incorrect data length. Expected %d got %d", 100, len(respData))
	}

	if !bytes.Equal(respData, sampleData[0:100]) {
		t.Fatalf("invalid data returned")
	}
}

func TestRangeGetSingleRelativeRangeStart(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(path.Join(dir, "data.txt"), sampleData, os.ModePerm)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)
	url := "http://" + listenAddress + "/data.txt"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", "rangetest/1.0")
	req.Header.Add("Range", "bytes=400-")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 206 {
		t.Fatalf("unexpected HTTP status code. Expected %d got %d", 206, resp.StatusCode)
	}

	if resp.ContentLength != int64(100) {
		t.Fatalf("incorrect value of content length header. Expected %d got %d", 100, resp.ContentLength)
	}

	ct, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		panic(err)
	}
	if ct != "text/plain" {
		t.Fatalf("incorrect value of content typ header. Expected 'text/plain' got '%s'", ct)
	}

	if value := resp.Header.Get("Content-Range"); value != "bytes 400-499/500" {
		t.Fatalf("incorrect or missing Accept-Ranges header. Expected 'bytes 400-499/500' got '%s'", value)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if len(respData) != 100 {
		t.Fatalf("incorrect data length. Expected %d got %d", 100, len(respData))
	}

	if !bytes.Equal(respData, sampleData[400:500]) {
		t.Fatalf("invalid data returned")
	}
}

func TestRangeGetSingleRelativeRangeEnd(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(path.Join(dir, "data.txt"), sampleData, os.ModePerm)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)
	url := "http://" + listenAddress + "/data.txt"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", "rangetest/1.0")
	req.Header.Add("Range", "bytes=-100")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 206 {
		t.Fatalf("unexpected HTTP status code. Expected %d got %d", 206, resp.StatusCode)
	}

	if resp.ContentLength != int64(100) {
		t.Fatalf("incorrect value of content length header. Expected %d got %d", 100, resp.ContentLength)
	}

	ct, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		panic(err)
	}
	if ct != "text/plain" {
		t.Fatalf("incorrect value of content typ header. Expected 'text/plain' got '%s'", ct)
	}

	if value := resp.Header.Get("Content-Range"); value != "bytes 400-499/500" {
		t.Fatalf("incorrect or missing Accept-Ranges header. Expected 'bytes 400-499/500' got '%s'", value)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if len(respData) != 100 {
		t.Fatalf("incorrect data length. Expected %d got %d", 100, len(respData))
	}

	if !bytes.Equal(respData, sampleData[400:500]) {
		t.Fatalf("invalid data returned")
	}
}

func TestRangeGetMultipleAbsoluteRanges(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(path.Join(dir, "data.txt"), sampleData, os.ModePerm)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)
	url := "http://" + listenAddress + "/data.txt"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", "rangetest/1.0")
	req.Header.Add("Range", "bytes=0-99,200-299,400-499")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 206 {
		t.Fatalf("unexpected HTTP status code. Expected %d got %d", 206, resp.StatusCode)
	}

	ct, args, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		panic(err)
	}
	if ct != "multipart/byteranges" {
		t.Fatalf("incorrect value of content typ header. Expected 'multipart/byteranges' got '%s'", ct)
	}
	boundary := args["boundary"]
	if boundary == "" {
		t.Fatalf("missing multipart boundary")
	}
	mpReader := multipart.NewReader(resp.Body, boundary)

	expectedContentRangeHeaders := []string{
		"bytes 0-99/500",
		"bytes 200-299/500",
		"bytes 400-499/500",
	}
	expectedData := [][]byte{
		sampleData[0:100],
		sampleData[200:300],
		sampleData[400:500],
	}
	expectedNumberOfParts := 3

	partIdx := 0
	for {
		part, err := mpReader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		if partIdx > expectedNumberOfParts-1 {
			t.Fatalf("unexpected number of data parts returned. Expected %d but got at least %d", expectedNumberOfParts, partIdx)
		}

		if partType := part.Header.Get("Content-Type"); partType != "text/plain" {
			t.Fatalf("invalid content type header value in part %d. Expected 'text/plain' got '%s'", partIdx+1, partType)
		}

		if contentRange := part.Header.Get("Content-Range"); contentRange != expectedContentRangeHeaders[partIdx] {
			t.Fatalf("invalid content range header value in part %d. Expected '%s' got '%s'", partIdx+1, expectedContentRangeHeaders[partIdx], contentRange)
		}

		partData, err := io.ReadAll(part)
		if err != nil {
			t.Fatalf("error reading data from part %d: %s", partIdx+1, err.Error())
		}

		if !bytes.Equal(partData, expectedData[partIdx]) {
			t.Fatalf("invalid data returned in part %d", partIdx+1)
		}

		partIdx++
	}
}

func TestRangeGetMultipleAbsoluteAndRelativeRanges1(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(path.Join(dir, "data.txt"), sampleData, os.ModePerm)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)
	url := "http://" + listenAddress + "/data.txt"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", "rangetest/1.0")
	req.Header.Add("Range", "bytes=0-99,-100")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 206 {
		t.Fatalf("unexpected HTTP status code. Expected %d got %d", 206, resp.StatusCode)
	}

	ct, args, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		panic(err)
	}
	if ct != "multipart/byteranges" {
		t.Fatalf("incorrect value of content typ header. Expected 'multipart/byteranges' got '%s'", ct)
	}
	boundary := args["boundary"]
	if boundary == "" {
		t.Fatalf("missing multipart boundary")
	}
	mpReader := multipart.NewReader(resp.Body, boundary)

	expectedContentRangeHeaders := []string{
		"bytes 0-99/500",
		"bytes 400-499/500",
	}
	expectedData := [][]byte{
		sampleData[0:100],
		sampleData[400:500],
	}
	expectedNumberOfParts := 2

	partIdx := 0
	for {
		part, err := mpReader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		if partIdx > expectedNumberOfParts-1 {
			t.Fatalf("unexpected number of data parts returned. Expected %d but got at least %d", expectedNumberOfParts, partIdx)
		}

		if partType := part.Header.Get("Content-Type"); partType != "text/plain" {
			t.Fatalf("invalid content type header value in part %d. Expected 'text/plain' got '%s'", partIdx+1, partType)
		}

		if contentRange := part.Header.Get("Content-Range"); contentRange != expectedContentRangeHeaders[partIdx] {
			t.Fatalf("invalid content range header value in part %d. Expected '%s' got '%s'", partIdx+1, expectedContentRangeHeaders[partIdx], contentRange)
		}

		partData, err := io.ReadAll(part)
		if err != nil {
			t.Fatalf("error reading data from part %d: %s", partIdx+1, err.Error())
		}

		if !bytes.Equal(partData, expectedData[partIdx]) {
			t.Fatalf("invalid data returned in part %d", partIdx+1)
		}

		partIdx++
	}
}

func TestRangeGetMultipleAbsoluteAndRelativeRanges2(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(path.Join(dir, "data.txt"), sampleData, os.ModePerm)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)
	url := "http://" + listenAddress + "/data.txt"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", "rangetest/1.0")
	req.Header.Add("Range", "bytes=0-99,400-")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 206 {
		t.Fatalf("unexpected HTTP status code. Expected %d got %d", 206, resp.StatusCode)
	}

	ct, args, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		panic(err)
	}
	if ct != "multipart/byteranges" {
		t.Fatalf("incorrect value of content typ header. Expected 'multipart/byteranges' got '%s'", ct)
	}
	boundary := args["boundary"]
	if boundary == "" {
		t.Fatalf("missing multipart boundary")
	}
	mpReader := multipart.NewReader(resp.Body, boundary)

	expectedContentRangeHeaders := []string{
		"bytes 0-99/500",
		"bytes 400-499/500",
	}
	expectedData := [][]byte{
		sampleData[0:100],
		sampleData[400:500],
	}
	expectedNumberOfParts := 2

	partIdx := 0
	for {
		part, err := mpReader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		if partIdx > expectedNumberOfParts-1 {
			t.Fatalf("unexpected number of data parts returned. Expected %d but got at least %d", expectedNumberOfParts, partIdx)
		}

		if partType := part.Header.Get("Content-Type"); partType != "text/plain" {
			t.Fatalf("invalid content type header value in part %d. Expected 'text/plain' got '%s'", partIdx+1, partType)
		}

		if contentRange := part.Header.Get("Content-Range"); contentRange != expectedContentRangeHeaders[partIdx] {
			t.Fatalf("invalid content range header value in part %d. Expected '%s' got '%s'", partIdx+1, expectedContentRangeHeaders[partIdx], contentRange)
		}

		partData, err := io.ReadAll(part)
		if err != nil {
			t.Fatalf("error reading data from part %d: %s", partIdx+1, err.Error())
		}

		if !bytes.Equal(partData, expectedData[partIdx]) {
			t.Fatalf("invalid data returned in part %d", partIdx+1)
		}

		partIdx++
	}
}

func TestRangeUnsupportedUnitType(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(path.Join(dir, "data.txt"), sampleData, os.ModePerm)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)
	url := "http://" + listenAddress + "/data.txt"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", "rangetest/1.0")
	req.Header.Add("Range", "centimeters=0-")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("unexpected HTTP status code. Expected %d got %d", 200, resp.StatusCode)
	}

	if resp.ContentLength != int64(len(sampleData)) {
		t.Fatalf("incorrect value of content length header. Expected %d got %d", len(sampleData), resp.ContentLength)
	}

	ct, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		panic(err)
	}
	if ct != "text/plain" {
		t.Fatalf("incorrect value of content typ header. Expected 'text/plain' got '%s'", ct)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if len(respData) != len(sampleData) {
		t.Fatalf("incorrect data length. Expected %d got %d", len(sampleData), len(respData))
	}

	if !bytes.Equal(respData, sampleData) {
		t.Fatalf("invalid data returned")
	}
}

func TestRangeStartIndexOutOfRange(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(path.Join(dir, "data.txt"), sampleData, os.ModePerm)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)
	url := "http://" + listenAddress + "/data.txt"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", "rangetest/1.0")
	req.Header.Add("Range", "bytes=700-800")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 416 {
		t.Fatalf("unexpected HTTP status code. Expected %d got %d", 416, resp.StatusCode)
	}
}

func TestRangeEndIndexOutOfRange(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(path.Join(dir, "data.txt"), sampleData, os.ModePerm)

	listenAddress := getListenAddress()

	server := router.New()
	server.ServeFiles(dir, "/")
	go func() {
		server.ListenAndServe(listenAddress)
	}()
	time.Sleep(5 * time.Millisecond)
	url := "http://" + listenAddress + "/data.txt"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", "rangetest/1.0")
	req.Header.Add("Range", "bytes=0-700")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 206 {
		t.Fatalf("unexpected HTTP status code. Expected %d got %d", 206, resp.StatusCode)
	}

	if resp.ContentLength != int64(len(sampleData)) {
		t.Fatalf("incorrect value of content length header. Expected %d got %d", len(sampleData), resp.ContentLength)
	}

	ct, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		panic(err)
	}
	if ct != "text/plain" {
		t.Fatalf("incorrect value of content typ header. Expected 'text/plain' got '%s'", ct)
	}

	if value := resp.Header.Get("Content-Range"); value != "bytes 0-499/500" {
		t.Fatalf("incorrect or missing Accept-Ranges header. Expected 'bytes 0-499/500' got '%s'", value)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if len(respData) != len(sampleData) {
		t.Fatalf("incorrect data length. Expected %d got %d", len(sampleData), len(respData))
	}

	if !bytes.Equal(respData, sampleData) {
		t.Fatalf("invalid data returned")
	}
}

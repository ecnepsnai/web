package router

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

// CacheMaxAge the amount of time browsers may consider static content to be fresh.
// Set this to 0 to not include a "Cache-Control" header for static requests.
var CacheMaxAge time.Duration = 24 * time.Hour

func (s *impl) serveStatic(dir, url string, w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" && req.Method != "HEAD" {
		s.MethodNotAllowedHandle(w, req)
		return
	}

	requestPath := stripPath(url)
	filePath := path.Join(dir, requestPath)

	s.log.PDebug("Serving static request", map[string]interface{}{
		"request_path": requestPath,
		"file_path":    filePath,
	})

	f, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		s.log.PInfo("Static file not found", map[string]interface{}{
			"request_path": requestPath,
			"file_path":    filePath,
		})
		s.NotFoundHandle(w, req)
		return
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		s.log.PError("Error getting static file info", map[string]interface{}{
			"request_path": requestPath,
			"file_path":    filePath,
			"error":        err.Error(),
		})
		s.NotFoundHandle(w, req)
		return
	}

	sendBody := req.Method == "GET"
	if modifiedSinceStr := req.Header.Get("If-Modified-Since"); modifiedSinceStr != "" {
		modifiedSince, err := httpDateToTime(modifiedSinceStr)
		if err != nil {
			modifiedSince = time.Now()
		}

		if info.ModTime().Sub(modifiedSince) < 0 {
			sendBody = false
		}
	}

	if CacheMaxAge > 0 {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d; public", int(CacheMaxAge.Seconds())))
	}
	w.Header().Set("Content-Type", MimeGetter.GetMime(filePath))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))
	w.Header().Add("Last-Modified", timeToHTTPDate(info.ModTime().UTC()))
	w.Header().Set("Date", timeToHTTPDate(time.Now().UTC()))
	if sendBody {
		io.Copy(w, f)
	} else {
		w.WriteHeader(204)
	}
}

const httpDateLayout = "Mon, 02 Jan 2006 15:04:05 GMT"

func httpDateToTime(date string) (time.Time, error) {
	return time.Parse(httpDateLayout, date)
}

func timeToHTTPDate(date time.Time) string {
	return date.Format(httpDateLayout)
}

func stripPath(inS string) (outS string) {
	// fast strip ../ from all paths
	in := []byte(inS)
	out := bytes.ReplaceAll(in, []byte{0x2e, 0x2e, 0x2f}, []byte{})
	outS = string(out)
	return
}

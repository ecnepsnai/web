package router

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// CacheMaxAge the amount of time browsers may consider static content to be fresh.
// Set this to 0 to not include a "Cache-Control" header for static requests.
var CacheMaxAge time.Duration = 24 * time.Hour

// IndexFileName is the name used when searching a directory for an index
var IndexFileName = "index.html"

// GenerateDirectoryListing if the router should generate a directory listing for static directories that do not have
// an index file (see also IndexFileName)
var GenerateDirectoryListing = true

func (s *impl) serveStatic(dir, url string, w http.ResponseWriter, req *http.Request) {
	requestPath := stripPath(url)
	shouldRenderDirectoryListing := false
	if requestPath == "" || strings.HasSuffix(requestPath, "/") {
		// First check if an index file is found
		if fileExists(path.Join(dir, requestPath+IndexFileName)) {
			requestPath += IndexFileName
		} else if fileExists(path.Join(dir, requestPath)) {
			// If an index file is not found, check if the directory exists
			shouldRenderDirectoryListing = true
		}
	}
	filePath := path.Join(dir, requestPath)

	if shouldRenderDirectoryListing {
		if !GenerateDirectoryListing {
			s.NotFoundHandle(w, req)
			return
		}

		s.makeDirectoryIndex(filePath, requestPath, w)
		return
	}

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

	if ranges := ParseRangeHeader(req.Header.Get("range")); len(ranges) > 0 && sendBody {
		headers := map[string]string{
			"Last-Modified": timeToHTTPDate(info.ModTime().UTC()),
		}
		if CacheMaxAge > 0 {
			headers["Cache-Control"] = fmt.Sprintf("max-age=%d; public", int(CacheMaxAge.Seconds()))
		}
		err = ServeHTTPRange(ServeHTTPRangeOptions{
			Headers:     headers,
			Ranges:      ranges,
			Reader:      f,
			TotalLength: uint64(info.Size()),
			MIMEType:    MimeGetter.GetMime(filePath),
			Writer:      w,
		})
		if err != nil {
			s.log.PError("Error serving ranged static file", map[string]interface{}{
				"request_path": requestPath,
				"file_path":    filePath,
				"error":        err.Error(),
			})
		}
		return
	}

	if CacheMaxAge > 0 {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d; public", int(CacheMaxAge.Seconds())))
	}
	w.Header().Set("Content-Type", MimeGetter.GetMime(filePath))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))
	w.Header().Add("Last-Modified", timeToHTTPDate(info.ModTime().UTC()))
	w.Header().Set("Date", timeToHTTPDate(time.Now().UTC()))
	w.Header().Set("Accept-Ranges", "bytes")
	if sendBody {
		io.Copy(w, f)
	} else {
		w.WriteHeader(200)
	}
}

// ServeHTTPRangeOptions options for serving a HTTP range request
type ServeHTTPRangeOptions struct {
	// Any additional headers to append to the request.
	// Do not specify a content-type here, instead use the MIMEType property.
	Headers map[string]string
	// Cookies to set on the response
	Cookies []http.Cookie
	// Byte ranges from the HTTP request
	Ranges []ByteRange
	// The incoming reader, must support seeking
	Reader io.ReadSeeker
	// The total length of the data
	TotalLength uint64
	// The content type of the data
	MIMEType string
	// The outgoing HTTP response writer
	Writer http.ResponseWriter
}

// ServeHTTPRange serve a HTTP range
func ServeHTTPRange(options ServeHTTPRangeOptions) error {
	for i := 0; i < len(options.Ranges); i++ {
		r := options.Ranges[i]
		if r.Start >= int64(options.TotalLength) {
			options.Writer.WriteHeader(416)
			return nil
		}
		if r.End >= int64(options.TotalLength) {
			options.Ranges[i].End = -1
		}
	}

	if len(options.Ranges) == 1 {
		return serveHTTPRangeSingle(options)
	}

	return serveHTTPRangeMulti(options)
}

func handleRange(reader io.ReadSeeker, writer io.Writer, r ByteRange) error {
	if r.Start >= 0 {
		if _, err := reader.Seek(r.Start, 0); err != nil {
			return err
		}
		if r.End >= 0 {
			// bytes=100-200
			if _, err := io.CopyN(writer, reader, (r.End-r.Start)+1); err != nil {
				return err
			}
		} else {
			// bytes=100-
			if _, err := io.Copy(writer, reader); err != nil {
				return err
			}
		}
	} else {
		if _, err := reader.Seek(r.End-(r.End*2), 2); err != nil {
			return err
		}
		// bytes=-100
		if _, err := io.Copy(writer, reader); err != nil {
			return err
		}
	}

	return nil
}

func serveHTTPRangeSingle(options ServeHTTPRangeOptions) error {
	r := options.Ranges[0]

	if CacheMaxAge > 0 {
		options.Writer.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d; public", int(CacheMaxAge.Seconds())))
	}
	options.Writer.Header().Set("Content-Type", options.MIMEType)
	options.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", r.Length(options.TotalLength)))
	options.Writer.Header().Set("Content-Range", r.ContentRangeValue(options.TotalLength))
	options.Writer.Header().Set("Date", timeToHTTPDate(time.Now().UTC()))
	options.Writer.WriteHeader(206)

	return handleRange(options.Reader, options.Writer, r)
}

func serveHTTPRangeMulti(options ServeHTTPRangeOptions) error {
	mp := multipart.NewWriter(options.Writer)
	options.Writer.Header().Set("Content-Type", fmt.Sprintf("multipart/byteranges; boundary=%s", mp.Boundary()))
	for k, v := range options.Headers {
		options.Writer.Header().Add(k, v)
	}
	for _, cookie := range options.Cookies {
		http.SetCookie(options.Writer, &cookie)
	}
	options.Writer.Header().Set("Date", timeToHTTPDate(time.Now().UTC()))
	options.Writer.WriteHeader(206)

	for _, r := range options.Ranges {
		part, err := mp.CreatePart(map[string][]string{
			"Content-Type":  {options.MIMEType},
			"Content-Range": {r.ContentRangeValue(options.TotalLength)},
		})
		if err != nil {
			return err
		}

		if err := handleRange(options.Reader, part, r); err != nil {
			return err
		}
	}

	mp.Close()
	return nil
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

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// ByteRange describes a range of offsets for reading from a byte slice.
//
// There are thee possabilities for byte ranges:
// (Start=100, End=200) Read data from offset 100 until offset 200.
// (Start=100, End=-1) Read all remaining data from offset 100 until the end of the reader.
// (Start=-1, End=100) Read the last 100 bytes of the reader.
type ByteRange struct {
	Start int64
	End   int64
}

// Length return the length of data represented by this byte range
func (br ByteRange) Length(total uint64) uint64 {
	start := uint64(0)
	end := total

	// bytes=-100
	if br.Start < 0 && br.End >= 0 {
		start = total - uint64(br.End)
		end = total
	}

	// bytes=100-200
	if br.Start >= 0 && br.End >= 0 {
		start = uint64(br.Start)
		end = uint64(br.End + 1)
	}

	// bytes=100-
	if br.Start >= 0 && br.End < 0 {
		start = uint64(br.Start)
		end = total
	}

	return end - start
}

// ContentRangeValue will return a value sutible for the content-range header
func (br ByteRange) ContentRangeValue(total uint64) string {
	start := ""
	end := ""

	// bytes=-100
	if br.Start < 0 && br.End >= 0 {
		start = fmt.Sprintf("%d", total-uint64(br.End))
		end = fmt.Sprintf("%d", total-1)
	}

	// bytes=100-200
	if br.Start >= 0 && br.End >= 0 {
		start = fmt.Sprintf("%d", br.Start)
		end = fmt.Sprintf("%d", br.End)
	}

	// bytes=100-
	if br.Start >= 0 && br.End < 0 {
		start = fmt.Sprintf("%d", br.Start)
		end = fmt.Sprintf("%d", total-1)
	}

	return fmt.Sprintf("bytes %s-%s/%d", start, end, total)
}

// ParseRangeHeader will parse the value from the HTTP ranges header and return a slice of byte ranges, or nil if the
// headers value is malformed.
func ParseRangeHeader(value string) []ByteRange {
	value = strings.ToLower(value)

	if len(value) < 5 {
		return nil
	}

	prefix := value[0:6]
	if prefix != "bytes=" {
		return nil
	}

	rangeStr := strings.ReplaceAll(value[6:], ", ", ",")
	ranges := strings.Split(rangeStr, ",")

	byteRanges := make([]ByteRange, len(ranges))
	for i, r := range ranges {
		parts := strings.Split(r, "-")
		if len(parts) < 1 {
			return nil
		}

		start := int64(-1)
		end := int64(-1)

		if parts[0] != "" {
			s, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return nil
			}
			start = s
		}
		if len(parts) == 2 && parts[1] != "" {
			e, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return nil
			}
			end = e
		}

		byteRanges[i] = ByteRange{start, end}
	}

	return byteRanges
}

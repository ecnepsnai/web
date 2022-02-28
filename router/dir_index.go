package router

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ecnepsnai/logtic"
)

//go:embed dir_index.html
var dirIndex string

type dirIndexTemplateType struct {
	Title       string
	Directories []string
	Files       []dirIndexFileType
	IsEmpty     bool
}

type dirIndexFileType struct {
	Name string
	Size string
}

func (s *impl) makeDirectoryIndex(dir, requestPath string, w http.ResponseWriter, req *http.Request) {
	s.log.PDebug("Serving directory listing", map[string]interface{}{
		"request_path":   requestPath,
		"directory_path": dir,
	})

	title := requestPath
	if title == "" {
		title = "/"
	}

	templateData := dirIndexTemplateType{
		Title: title,
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		s.log.PError("Error reading directory", map[string]interface{}{
			"dir":   dir,
			"error": err.Error(),
		})
		w.WriteHeader(500)
		return
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		if entry.IsDir() {
			templateData.Directories = append(templateData.Directories, entry.Name()+"/")
		} else {
			templateData.Files = append(templateData.Files, dirIndexFileType{
				Name: entry.Name(),
				Size: logtic.FormatBytesB(uint64(info.Size())),
			})
		}
	}
	if len(entries) == 0 {
		templateData.IsEmpty = true
	}

	t, err := template.New("index").Parse(dirIndex)
	if err != nil {
		s.log.PError("Error forming template for directory index", map[string]interface{}{
			"error": err.Error(),
		})
		w.WriteHeader(500)
		return
	}

	buf := &bytes.Buffer{}
	if err := t.ExecuteTemplate(buf, "main", templateData); err != nil {
		s.log.PError("Error executing template for directory index", map[string]interface{}{
			"error": err.Error(),
		})
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	w.Header().Set("Date", timeToHTTPDate(time.Now().UTC()))
	io.Copy(w, buf)
}

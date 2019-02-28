package zpages

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/awbraunstein/zpages/buildinfo"
)

type dynamicInfo struct {
	NumGoroutine int
}

type templateInfo struct {
	Goarch      string
	Goos        string
	Gc          string
	GoVersion   string
	Environment []string
	DynamicData dynamicInfo
	BuildTime   string
	CommitHash  string
	Hostname    string
}

type Statusz struct {
	tmpl *template.Template
	data templateInfo
}

func NewStatusz() (*Statusz, error) {
	tmplPath := filepath.Join(os.Getenv("GOPATH"), "src/github.com/awbraunstein/zpages", "templates/statusz.tmpl")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return nil, err
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Could not read hostname; err=%v", err)
	}
	return &Statusz{
		tmpl: tmpl,
		data: templateInfo{
			Goarch:      runtime.GOARCH,
			Goos:        runtime.GOOS,
			Gc:          runtime.Compiler,
			GoVersion:   runtime.Version(),
			Environment: os.Environ(),
			BuildTime:   buildinfo.BuildTime,
			CommitHash:  buildinfo.CommitHash,
			Hostname:    hostname,
		},
	}, nil
}

func (h *Statusz) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	h.data.DynamicData.NumGoroutine = runtime.NumGoroutine()
	if err := h.tmpl.Execute(resp, h.data); err != nil {
		log.Println("Unable to execute statusz.tmpl; err=%v", err)
	}
}

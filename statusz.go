package zpages

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/awbraunstein/zpages/buildinfo"
)

// When this process was started. Not entirely accurate but probably close
// enough.
// TODO: Get this information from the OS.
var startTime = time.Now()

type dynamicInfo struct {
	NumGoroutine int
}

// TODO: Allow users to add custom sub-templates to be rendered.
type templateInfo struct {
	// Information determined at creation of the handler.
	Invocation string
	Pid        int
	LaunchTime string
	Goarch     string
	Goos       string
	Gc         string
	GoVersion  string
	BuildTime  string
	CommitHash string
	Hostname   string

	// Information when rendering the page.
	Runtime      string
	Environment  []string
	NumGoroutine int
	OtherData    map[string]string
}

// Statusz is a handler that renders the statusz page.
type Statusz struct {
	tmpl         *template.Template
	data         templateInfo
	dataAdderFns []func(map[string]string)
}

// NewStatusz creates a new Statusz handler. Passed in functions can add
// additional data to the rendered template.
func NewStatusz(dataAdderFns ...func(map[string]string)) (*Statusz, error) {
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
		tmpl:         tmpl,
		dataAdderFns: dataAdderFns,
		data: templateInfo{
			Invocation: strings.Join(os.Args, " "),
			Pid:        os.Getpid(),
			LaunchTime: startTime.Format(time.RFC822),
			Goarch:     runtime.GOARCH,
			Goos:       runtime.GOOS,
			Gc:         runtime.Compiler,
			GoVersion:  runtime.Version(),
			BuildTime:  buildinfo.BuildTime,
			CommitHash: buildinfo.CommitHash,
			Hostname:   hostname,
		},
	}, nil
}

func (h *Statusz) ServeHTTP(resp http.ResponseWriter, _ *http.Request) {
	h.data.NumGoroutine = runtime.NumGoroutine()
	h.data.Environment = os.Environ()
	h.data.Runtime = time.Since(startTime).String()
	// Clear the OtherData map each time.
	h.data.OtherData = make(map[string]string)
	for _, fn := range h.dataAdderFns {
		fn(h.data.OtherData)
	}
	if err := h.tmpl.Execute(resp, h.data); err != nil {
		log.Printf("Unable to execute statusz.tmpl; err=%v", err)
	}
}

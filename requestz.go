package zpages

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Allow for easy mocking of time.Now()
var timeNow = time.Now

const requestRetention = time.Duration(2) * time.Minute

// RequestInfo holds the information of a single request.
type RequestInfo struct {
	Timestamp time.Time
	Status    int
	Request   *http.Request
}

// Requestz is an http handler that renders the requestz page.
type Requestz struct {
	tmpl              *template.Template
	completedRequests []*RequestInfo
	additionalDataFns []func(*RequestInfo)
	mu                sync.RWMutex
}

// Create a new Requestz handler. Additional functions can modify the current
// request's RequestInfo and are run in order.
func NewRequestz(additionalDataFns ...func(*RequestInfo)) (*Requestz, error) {
	tmplPath := filepath.Join(os.Getenv("GOPATH"), "src/github.com/awbraunstein/zpages", "templates/requestz.tmpl")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return nil, err
	}
	return &Requestz{
		tmpl:              tmpl,
		additionalDataFns: additionalDataFns,
	}, nil
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) Write(b []byte) (int, error) {
	if rec.status == 0 {
		rec.WriteHeader(200)
	}
	return rec.ResponseWriter.Write(b)
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func (h *Requestz) addRequest(status int, r *http.Request) {
	h.mu.Lock()
	// First try to delete any requests outside of the request retention.
	foundInBoundRequest := false
	for i, reqInfo := range h.completedRequests {
		if timeNow().Sub(reqInfo.Timestamp) <= requestRetention {
			h.completedRequests = h.completedRequests[i:]
			foundInBoundRequest = true
			break
		}
	}
	if !foundInBoundRequest {
		// This means that all the completed requests should be tossed.
		h.completedRequests = nil
	}
	ri := &RequestInfo{
		Timestamp: timeNow(),
		Status:    status,
		Request:   r,
	}
	for _, fn := range h.additionalDataFns {
		fn(ri)
	}
	// Then add the request.
	h.completedRequests = append(h.completedRequests, ri)
	h.mu.Unlock()
}

// Middleware allows for easy chaning of the Middleware handler.
func (h *Requestz) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sr := statusRecorder{ResponseWriter: w}
		next.ServeHTTP(&sr, r)
		// We don't want to block the request from returning so we are
		// doing this in a goroutine.
		go h.addRequest(sr.status, r)
	})
}

// Holds the data for requests for a given path.
type pathData struct {
	TimeBuckets map[time.Duration][]*RequestInfo
	Errors      []*RequestInfo
}

// requestzTmplData holds the data that is required to render the requestz
// template.
type requestzTmplData struct {
	RequestsByPath map[string]*pathData
}

func (requestzTmplData) GetDurationString(d time.Duration) string {
	return fmt.Sprintf("<%ds", d/time.Second)
}

func (requestzTmplData) GetElementId(path string, d time.Duration, tag string) string {
	normalizedPath := strings.ReplaceAll(path, " ", "-")
	return fmt.Sprintf("%s-%d-%s", normalizedPath, d/time.Second, tag)
}

// The time bucketing we will be using. These are in seconds and represent that
// are younger than that many seconds. A request may be in multiple buckets.
var timeBuckets = []int{10, 30, 60, 90, 120}

func (h *Requestz) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	tmplData := requestzTmplData{
		RequestsByPath: make(map[string]*pathData),
	}

	h.mu.RLock()
	for _, reqInfo := range h.completedRequests {
		// The path will be "METHOD path".
		path := fmt.Sprintf("%s %s", reqInfo.Request.Method, reqInfo.Request.URL.Path)
		pd := tmplData.RequestsByPath[path]
		if pd == nil {
			tmplData.RequestsByPath[path] = &pathData{
				TimeBuckets: make(map[time.Duration][]*RequestInfo),
			}
			pd = tmplData.RequestsByPath[path]
		}
		for _, bucket := range timeBuckets {
			tb := time.Duration(bucket) * time.Second
			// Make sure all the time buckets have an entry, even if
			// there are no requests in it.
			if pd.TimeBuckets[tb] == nil {
				pd.TimeBuckets[tb] = []*RequestInfo{}
			}
			if timeNow().Sub(reqInfo.Timestamp) < tb {
				pd.TimeBuckets[tb] = append(pd.TimeBuckets[tb], reqInfo)
			}
		}
		if reqInfo.Status != http.StatusOK {
			pd.Errors = append(pd.Errors, reqInfo)
		}
	}
	h.mu.RUnlock()
	if err := h.tmpl.Execute(resp, tmplData); err != nil {
		log.Printf("Unable to execute requestz.tmpl; err=%v", err)
	}
}

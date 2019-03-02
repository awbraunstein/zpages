package zpages

import "net/http"

// Healthz renders a page that displays a health check value. By default, this
// handler will render ok when the server is ready to serve. It will echo the
// healthString url param if it is provided. Use NewHealthz() to construct this
// struct.
type Healthz struct {
	okValue string
}

// NewHealthz creates a new Healthz http handler. This method accepts either 0
// or 1 params and will panic on any other number of values.
//
//  NewHealthz()
// or
//  NewHealthz("ok")
func NewHealthz(okValue ...string) *Healthz {
	if len(okValue) > 1 {
		panic("NewHealthz only accepts 0 or 1 arguments")
	}
	health := "ok"
	if len(okValue) == 1 {
		health = okValue[0]
	}
	return &Healthz{okValue: health}
}

func (h *Healthz) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if healthString, ok := req.URL.Query()["healthString"]; ok {
		resp.Write([]byte(healthString[0]))
		return
	}
	resp.Write([]byte(h.okValue))
}

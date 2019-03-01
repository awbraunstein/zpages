package zpages

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStatusz(t *testing.T) {
	handler, err := NewStatusz()
	if err != nil {
		t.Fatalf("unable to create new statusz handler; err=%v", err)
	}
	ts := httptest.NewServer(handler)
	defer ts.Close()
	client := ts.Client()
	res, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("Unable to get url:%q; err=%v", ts.URL, err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatalf("Unable to read response; err=%v", err)
	}
	if !strings.Contains(string(body), "Statusz") {
		t.Errorf(`Expected response to contain "Statusz" but didn't`)
	}
}

package zpages

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func sendRequest(t *testing.T, c *http.Client, url string) {
	resp, err := c.Get(url)
	if err != nil {
		t.Fatalf("unable to get url:%q; err=%v", url, err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("response status is %d, but expected 200", resp.StatusCode)
	}
}

func TestRecordsRequest(t *testing.T) {
	defer func(old func() time.Time) {
		timeNow = old
	}(timeNow)
	// Stop time.
	timeNow = func() time.Time {
		return time.Unix(1, 0)
	}

	done := make(chan bool)
	requestz, err := NewRequestz(func(_ *RequestInfo) {
		done <- true
	})
	if err != nil {
		t.Fatalf("could not create new requestz; err=%v", err)
	}
	ts := httptest.NewServer(requestz.Middleware(requestz))
	defer ts.Close()
	client := ts.Client()
	sendRequest(t, client, ts.URL)
	<-done
	// Request a lock to ensure the completedRequests slice has been updated.
	requestz.mu.RLock()
	if got := len(requestz.completedRequests); got != 1 {
		t.Fatalf("expected 1 request, but got %d", got)
	}
	requestz.mu.RUnlock()
	sendRequest(t, client, ts.URL)
	<-done
	// Request a lock to ensure the completedRequests slice has been updated.
	requestz.mu.RLock()
	if got := len(requestz.completedRequests); got != 2 {
		t.Fatalf("expected 2 request, but got %d", got)
	}
	requestz.mu.RUnlock()
}

func TestRemovesOldRequests(t *testing.T) {
	defer func(old func() time.Time) {
		timeNow = old
	}(timeNow)

	// Stop time.
	timeNow = func() time.Time {
		return time.Unix(1, 0)
	}

	done := make(chan bool)
	requestz, err := NewRequestz(func(_ *RequestInfo) {
		done <- true
	})
	if err != nil {
		t.Fatalf("could not create new requestz; err=%v", err)
	}
	ts := httptest.NewServer(requestz.Middleware(requestz))
	defer ts.Close()
	client := ts.Client()
	sendRequest(t, client, ts.URL)
	<-done
	requestz.mu.RLock()
	if got := len(requestz.completedRequests); got != 1 {
		t.Fatalf("expected 1 request, but got %d", got)
	}
	requestz.mu.RUnlock()
	// Move time to the FUTURE!
	timeNow = func() time.Time {
		return time.Unix(1000, 0)
	}
	sendRequest(t, client, ts.URL)
	<-done
	requestz.mu.RLock()
	if got := len(requestz.completedRequests); got != 1 {
		t.Fatalf("expected 1 request, but got %d", got)
	}
	requestz.mu.RUnlock()
	// Move time to 10s the FUTURE!
	timeNow = func() time.Time {
		return time.Unix(1010, 0)
	}
	sendRequest(t, client, ts.URL)
	<-done
	requestz.mu.RLock()
	if got := len(requestz.completedRequests); got != 2 {
		t.Fatalf("expected 2 request, but got %d", got)
	}
	requestz.mu.RUnlock()
	// Move time to requestRentention-1s the FUTURE!
	timeNow = func() time.Time {
		return time.Unix(1010-1+int64(requestRetention/time.Second), 0)
	}
	sendRequest(t, client, ts.URL)
	<-done
	requestz.mu.RLock()
	if got := len(requestz.completedRequests); got != 2 {
		t.Fatalf("expected 2 request, but got %d", got)
	}
	requestz.mu.RUnlock()
}

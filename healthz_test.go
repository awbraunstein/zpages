package zpages

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	tests := []struct {
		name         string
		defaultValue string
		param        string
		want         string
	}{
		{
			name:         "No param expects OK",
			defaultValue: "",
			param:        "",
			want:         "ok",
		}, {
			name:         "No param expects newDefault",
			defaultValue: "newDefault",
			param:        "",
			want:         "newDefault",
		}, {
			name:         "Expects echo param",
			defaultValue: "",
			param:        "echo",
			want:         "echo",
		}, {
			name:         "Expects foo param",
			defaultValue: "",
			param:        "foo",
			want:         "foo",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			handler := NewHealthz()
			if tc.defaultValue != "" {
				handler = NewHealthz(tc.defaultValue)
			}
			ts := httptest.NewServer(handler)
			defer ts.Close()
			client := ts.Client()
			url := ts.URL
			if tc.param != "" {
				url = url + "?healthString=" + tc.param
			}
			res, err := client.Get(url)
			if err != nil {
				t.Fatalf("Unable to get url:%q; err=%v", url, err)
			}
			got, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				t.Fatalf("Unable to read response; err=%v", err)
			}
			if got := string(got); got != tc.want {
				t.Errorf("Expected: %q, but got: %q", tc.want, got)
			}
		})
	}
}

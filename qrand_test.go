package qrand_test

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bitfield/qrand"
	"github.com/google/go-cmp/cmp"
)

func TestBytes(t *testing.T) {
	t.Parallel()
	called := false
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		wantURL := "/API/jsonI.php?length=64&type=uint8"
		if !cmp.Equal(wantURL, r.URL.String()) {
			t.Error(cmp.Diff(wantURL, r.URL.String()))
		}
		w.WriteHeader(http.StatusOK)
		data, err := os.Open("testdata/response.json")
		if err != nil {
			t.Fatal(err)
		}
		defer data.Close()
		io.Copy(w, data)
	}))
	defer ts.Close()
	qrand.HTTPClient = ts.Client()
	qrand.URL = ts.URL
	var got = make([]byte, 64)
	bytesRead, err := qrand.Read(got)
	if err != nil {
		t.Fatal(err)
	}
	if bytesRead != 64 {
		t.Errorf("want 64 bytes read, got %d", bytesRead)
	}
	if got[0] != 63 {
		t.Errorf("first byte should be 63, but was %d", got[0])
	}
	if !called {
		t.Error("want API call, but got none")
	}
}

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()
	jData, err := ioutil.ReadFile("testdata/response.json")
	if err != nil {
		t.Fatal(err)
	}
	want := qrand.APIResponse{[]byte{63, 25, 21, 239, 178, 131, 81, 166, 228, 3, 236, 255, 228, 216, 52, 195, 139, 187, 223, 161, 45, 5, 175, 173, 47, 57, 24, 212, 196, 34, 195, 132, 224, 211, 212, 6, 85, 135, 159, 57, 155, 213, 23, 33, 239, 85, 255, 76, 163, 51, 15, 251, 45, 216, 100, 123, 171, 8, 209, 92, 220, 207, 73, 172}}
	got := qrand.APIResponse{}
	err = json.Unmarshal(jData, &got)
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

type zeroReader struct{}

func (z zeroReader) Read(b []byte) (int, error) {
	copy(b, make([]byte, len(b)))

	return len(b), nil
}

func TestReassignReader(t *testing.T) {
	// not concurrency-safe, because we mess with global Reader
	got := make([]byte, 3)
	want := make([]byte, 3)
	origReader := qrand.Reader
	qrand.Reader = zeroReader{}
	bytesRead, err := qrand.Read(got)
	qrand.Reader = origReader
	if err != nil {
		t.Fatal(err)
	}
	if bytesRead != 3 {
		t.Errorf("want 3 bytes read, got %d", bytesRead)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

package qrand_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitfield/qrand"
	"github.com/google/go-cmp/cmp"
)

// func TestReadCannedData(t *testing.T) {
// 	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		data, err := os.Open("testdata/response.json")
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		w.WriteHeader(http.StatusOK)
// 		defer data.Close()
// 		io.Copy(w, data)
// 	}))
// 	defer ts.Close()
// 	qrand.HTTPClient = ts.Client()
// 	qrand.URL = ts.URL
// 	const numBytes = 8
// 	got := make([]byte, numBytes)
// 	n, err := qrand.Read(got)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if n != numBytes {
// 		t.Errorf("want %d bytes, got %d", numBytes, n)
// 	}
// }

func TestBytes(t *testing.T) {
	called := false
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		wantURL := "/API/jsonI.php?length=1&type=uint8&size=8"
		if !cmp.Equal(wantURL, r.URL.String()) {
			t.Error(cmp.Diff(wantURL, r.URL.String()))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	qrand.HTTPClient = ts.Client()
	qrand.URL = ts.URL
	var got bytes.Buffer
	err := qrand.Bytes(&got, 8)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("want API call, but got none")
	}
}

type ZeroReader struct{}

func (z ZeroReader) Read(b []byte) (n int, err error) {
	for i := range b {
		b[i] = 0
	}
	return len(b), nil
}

func TestReadZeroReader(t *testing.T) {
	qrand.Reader = ZeroReader{}
	const numBytes = 8
	got := make([]byte, numBytes)
	n, err := qrand.Read(got)
	if err != nil {
		t.Fatal(err)
	}
	if n != numBytes {
		t.Errorf("want %d bytes, got %d", numBytes, n)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	jData, err := ioutil.ReadFile("testdata/response.json")
	if err != nil {
		t.Fatal(err)
	}
	want := qrand.APIResponse{[]byte("5258aa2852307702")}
	got := qrand.APIResponse{}
	err = json.Unmarshal(jData, &got)
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(got, want) {
		t.Error(cmp.Diff(got, want))
	}
}

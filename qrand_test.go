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
		wantURL := "/API/jsonI.php?length=1&type=uint8&size=12"
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
	var got = make([]byte, 12)
	bytesRead, err := qrand.Read(got)
	if err != nil {
		t.Fatal(err)
	}
	if bytesRead != 12 {
		t.Errorf("want 12 bytes read, got %d", bytesRead)
	}
	if !called {
		t.Error("want API call, but got none")
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

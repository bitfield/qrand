package qrand_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bitfield/qrand"
	"github.com/google/go-cmp/cmp"
)

func TestReadMakesCorrectAPIRequestAndParsesResult(t *testing.T) {
	t.Parallel()
	called := false
	wantAPIKey := "dummyKey"
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		gotAPIKey := r.Header.Get("x-api-key")
		if wantAPIKey != gotAPIKey {
			t.Error("bad x-api-key header", cmp.Diff(wantAPIKey, gotAPIKey))
		}
		wantURL := "/?length=64&type=uint8"
		if wantURL != r.URL.String() {
			t.Error("bad request URI", cmp.Diff(wantURL, r.URL.String()))
		}
		http.ServeFile(w, r, "testdata/response.json")
	}))
	defer ts.Close()
	q := qrand.NewReader(wantAPIKey)
	q.HTTPClient = ts.Client()
	q.BaseURL = ts.URL
	got := make([]byte, 64)
	bytesRead, err := q.Read(got)
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

func TestReadReturnsErrorsWhenExpected(t *testing.T) {
	t.Parallel()
	q := qrand.NewReader("dummyKey")
	_, err := q.Read(make([]byte, 1025))
	if err == nil {
		t.Fatal("want error when requested bytes exceeds API maximum")
	}
	q.BaseURL = string([]rune{0x7F}) // invalid character in URLs
	_, err = q.Read([]byte{})
	if err == nil {
		t.Fatal("want error when base URL unparsable")
	}
	q.BaseURL = "invalid"
	_, err = q.Read([]byte{})
	if err == nil {
		t.Fatal("want error when base URL unreachable")
	}
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	q.HTTPClient = ts.Client()
	q.BaseURL = ts.URL
	_, err = q.Read([]byte{})
	if err == nil {
		t.Fatal("want error when API returns forbidden status")
	}
	ts = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1") // causes unexpected EOF on read
	}))
	q.HTTPClient = ts.Client()
	q.BaseURL = ts.URL
	_, err = q.Read([]byte{})
	if err == nil {
		t.Fatal("want error when response body can't be read")
	}
	ts = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	q.HTTPClient = ts.Client()
	q.BaseURL = ts.URL
	_, err = q.Read([]byte{})
	if err == nil {
		t.Fatal("want error on unexpected HTTP response status")
	}
	ts = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "") // invalid JSON
	}))
	q.HTTPClient = ts.Client()
	q.BaseURL = ts.URL
	_, err = q.Read([]byte{})
	if err == nil {
		t.Fatal("want error on invalid JSON response")
	}
}

func TestUnmarshalJSON_UnmarshalsValidData(t *testing.T) {
	t.Parallel()
	jData, err := os.ReadFile("testdata/response.json")
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

func TestUnmarshalJSON_ReturnsErrorsWhenExpected(t *testing.T) {
	t.Parallel()
	r := qrand.APIResponse{}
	err := r.UnmarshalJSON([]byte{})
	if err == nil {
		t.Fatal("want error for invalid JSON")
	}
	err = r.UnmarshalJSON([]byte(`{"nodata":"foryou"}`))
	if err == nil {
		t.Fatal("want error for missing 'data' field")
	}
	err = r.UnmarshalJSON([]byte(`{"data":1}`))
	if err == nil {
		t.Fatal("want error when 'data' is not an array")
	}
	err = r.UnmarshalJSON([]byte(`{"data":[]}`))
	if err == nil {
		t.Fatal("want error when 'data' has no elements")
	}
	err = r.UnmarshalJSON([]byte(`{"data":["not a number"]}`))
	if err == nil {
		t.Fatal("want error when 'data' has non-numeric element")
	}
	err = r.UnmarshalJSON([]byte(`{"data":[256]}`))
	if err == nil {
		t.Fatal("want error when 'data' element is too big for a byte")
	}
}

func TestSourceReturnsExpectedResultFromIntn(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/response.json")
	}))
	defer ts.Close()
	q := qrand.NewReader("dummyKey")
	q.HTTPClient = ts.Client()
	q.BaseURL = ts.URL
	rnd := rand.New(qrand.NewSource(q))
	want := 3
	got := rnd.Intn(10)
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

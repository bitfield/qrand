// Package qrand provides random numbers using the hardware quantum random
// number generator at the Australian National University (ANU). See
// https://quantumnumbers.anu.edu.au/ for details of this API.
package qrand

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// A qReader represents a client for the AQN API.
type qReader struct {
	// BaseURL holds the 'SCHEME://HOST' part of the request URL. This can
	// be overridden (for example, for testing against a local HTTP
	// server).
	BaseURL string
	// HTTPClient holds the *[http.Client] that will be used for requests.
	HTTPClient *http.Client
	apiKey     string
}

// NewReader creates and returns a qReader struct representing an AQN API
// client.
func NewReader(apiKey string) *qReader {
	return &qReader{
		BaseURL: "https://api.quantumnumbers.anu.edu.au",
		apiKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Read attempts to read enough random data from the API to fill buf, returning
// the number of bytes successfully read, along with any error.
func (q qReader) Read(buf []byte) (n int, err error) {
	if len(buf) > 1024 {
		return 0, fmt.Errorf("number of bytes must be less than 1024 (API limit): %d", len(buf))
	}
	size := len(buf)
	URL := fmt.Sprintf("%s?length=%d&type=uint8", q.BaseURL, size)
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("x-api-key", q.apiKey)
	resp, err := q.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusForbidden {
		return 0, errors.New("unauthorised: check your API key is valid")
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected response status %q", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("reading response body: %w", err)
	}
	r := APIResponse{}
	err = json.Unmarshal(data, &r)
	if err != nil {
		return 0, fmt.Errorf("invalid response: %w", err)
	}
	return copy(buf, r.Data), nil
}

// Source is a randomness source which implements rand.Source.
type source struct {
	Reader io.Reader
}

// Seed is a no-op, because a qrand source doesn't need seeding.
func (s *source) Seed(seed int64) {}

// Uint64 returns a random 64-bit value as a uint64.
func (s *source) Uint64() (value uint64) {
	binary.Read(s.Reader, binary.BigEndian, &value)
	return value
}

// Int63 returns a non-negative 63-bit integer as an int64.
func (s *source) Int63() (value int64) {
	return int64(s.Uint64() & ^uint64(1<<63))
}

// NewSource creates a [rand.Source] using q as the reader.
func NewSource(q *qReader) rand.Source {
	return &source{
		Reader: q,
	}
}

// APIResponse represents a response from the ANU QRNG API.
type APIResponse struct {
	Data []byte
}

// UnmarshalJSON reads the byte data in the raw API response into the
// APIResponse's Data field.
func (r *APIResponse) UnmarshalJSON(input []byte) error {
	raw := map[string]interface{}{}
	err := json.Unmarshal(input, &raw)
	if err != nil {
		return err
	}
	resp, ok := raw["data"]
	if !ok {
		return fmt.Errorf("no 'data' field found in response: %v", raw)
	}
	data, ok := resp.([]interface{})
	if !ok {
		return fmt.Errorf("want []interface{} value for 'data', got %T: %q", raw["data"], raw["data"])
	}
	if len(data) == 0 {
		return fmt.Errorf("not enough 'data' elements in response: %v", raw)
	}
	var rawVal float64
	for _, b := range data {
		if rawVal, ok = b.(float64); !ok {
			return fmt.Errorf("element '%v' in data should be a float64, but is a %T", b, b)
		}
		if rawVal > 255 {
			return fmt.Errorf("element '%f' is too big for a byte", rawVal)
		}
		r.Data = append(r.Data, byte(rawVal))
	}
	return nil
}

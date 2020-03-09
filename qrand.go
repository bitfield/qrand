// Package qrand provides random numbers using the hardware quantum random
// number generator at the Australian National University (ANU). See
// https://qrng.anu.edu.au/API/api-demo.php for details of the ANU QRNG API.
package qrand

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const maxBytesPerRequest = 1024 // API limit

var (
	// HTTPClient is the `*http.Client` which will be used to make API
	// requests. It has a 5-second timeout. To use a different timeout, set
	// HTTPClient.Timeout. To use a different client, set HTTPClient.
	HTTPClient *http.Client = &http.Client{
		Timeout: 5 * time.Second,
	}
	// URL is the URL of the ANU QRNG API server. To use a different server (for example for testing), set the URL accordingly.
	URL string = "https://qrng.anu.edu.au"
)

// Read calls the ANU QRNG API to read enough bytes to fill 'buf'. It returns
// the number of bytes actually read, or an error.
func Read(buf []byte) (n int, err error) {
	if len(buf) > maxBytesPerRequest {
		return 0, fmt.Errorf("number of bytes must be less than %d (API limit): %d", maxBytesPerRequest, len(buf))
	}
	size := len(buf)
	resp, err := HTTPClient.Get(fmt.Sprintf("%s/API/jsonI.php?length=%d&type=uint8", URL, size))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("reading response body: %w", err)
	}
	resp.Body.Close()
	respString := string(respBytes)
	resp.Body = ioutil.NopCloser(strings.NewReader(respString))
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected response status %d: %q", resp.StatusCode, respString)
	}
	var r = APIResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return 0, fmt.Errorf("decoding error for %q: %w", respString, err)
	}
	copy(buf, r.Data)
	return len(buf), nil
}

// APIResponse represents a response from the ANU QRNG API.
type APIResponse struct {
	Data []byte `json:"data"`
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
		return fmt.Errorf("No 'data' field found in response: %v", raw)
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

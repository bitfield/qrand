package qrand

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const max = 1024 * 1024 // API limit

var (
	Reader     io.Reader
	HTTPClient *http.Client = &http.Client{
		Timeout: 5 * time.Second,
	}
	URL string = "qrng.anu.edu.au"
)

func Bytes(buf *bytes.Buffer, n int) error {
	if n > max {
		return fmt.Errorf("number of bytes must be less than %d (API limit): %d", max, n)
	}
	blocks := 1
	size := n
	if n > 1024 {
		size = 1024
		blocks = n / 1024
		if n%1024 > 0 {
			blocks++
		}
	}
	HTTPClient.Get(fmt.Sprintf("%s/API/jsonI.php?length=%d&type=uint8&size=%d", URL, blocks, size))
	return nil
}

type APIResponse struct {
	Data []byte `json:"data"`
}

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
	value, ok := data[0].(string)
	if !ok {
		return fmt.Errorf("want string data value, got %T: %v", data[0], data[0])
	}
	r.Data = []byte(value)
	return nil
}

func Read(b []byte) (n int, err error) {
	return io.ReadFull(Reader, b)
}

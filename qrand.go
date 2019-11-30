package qrand

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	Reader     io.Reader
	HTTPClient *http.Client
	URL        string
)

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

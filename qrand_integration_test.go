//go:build integration
// +build integration

package qrand_test

import (
	"os"
	"testing"

	"github.com/bitfield/qrand"
)

func TestIntegration(t *testing.T) {
	apiKey := os.Getenv("AQN_API_KEY")
	if apiKey == "" {
		t.Fatal("AQN_API_KEY must be set")
	}
	data := make([]byte, 64)
	q := qrand.NewReader(apiKey)
	n, err := q.Read(data)
	if err != nil {
		t.Error(err)
	}
	if n != 64 {
		t.Errorf("want 64 bytes, got %d", n)
	}
}

// +build integration

package qrand_test

import (
	"testing"

	"github.com/bitfield/qrand"
)

func TestIntegration(t *testing.T) {
	data := make([]byte, 64)
	n, err := qrand.Read(data)
	if err != nil {
		t.Error(err)
	}
	if n != 64 {
		t.Errorf("want 64 bytes, got %d", n)
	}
}

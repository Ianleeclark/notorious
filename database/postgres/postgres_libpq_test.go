package postgres

import (
	"testing"
)

func TestNewPGListener(t *testing.T) {
	_, err := NewListener(CONFIG)
	if err != nil {
		t.Fatalf("Received err in TestNewPGListener: %v", err)
	}
}

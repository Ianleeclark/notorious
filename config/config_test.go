package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	loadedConfig := LoadConfig()

	if loadedConfig.MySQLHost != "localhost" {
		t.Fatalf("Expected %s, got %v", "localhost", loadedConfig.MySQLHost)
	}
}

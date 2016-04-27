package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
    loadedConfig := LoadConfig()

    if loadedConfig.MySQLHost != "127.0.0.1" {
        t.Fatalf("Expected %s, got %v", "127.0.0.1", loadedConfig.MySQLHost)
    }
}


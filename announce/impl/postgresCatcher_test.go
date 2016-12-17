package catcherImpl

import (
	"github.com/GrappigPanda/notorious/config"
	"testing"
)

var CONFIG = config.ConfigStruct{
	"postgres",
	"localhost",
	"5432",
	"postgres",
	"",
	"testdb",
	false,
}

func TestNewCatcher(t *testing.T) {
	_ = NewCatcher(CONFIG)
}

func TestHandleTorrent(t *testing.T) {
	catcher := NewCatcher(CONFIG)
	catcher.HandleNewTorrent()
}

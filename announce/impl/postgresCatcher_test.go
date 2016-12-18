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
	_ = NewPostgresCatcher(CONFIG)
}

func TestHandleTorrent(t *testing.T) {
	catcher := NewPostgresCatcher(CONFIG)
	catcher.HandleNewTorrent()
}

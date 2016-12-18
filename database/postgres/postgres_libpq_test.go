package postgres

import (
	"github.com/GrappigPanda/notorious/config"
	"github.com/GrappigPanda/notorious/database/schemas"
	"github.com/lib/pq"
	"os"
	"testing"
	"time"
)

var CONFIG = config.ConfigStruct{
	"postgres",
	"localhost",
	"5432",
	"postgres",
	"",
	"testdb",
	false,
	nil,
}

var LISTENER, ERR = NewListener(CONFIG)
var CALLBACKGLOBAL = ""

func TestNewPGListener(t *testing.T) {
	if ERR != nil {
		t.Fatalf("Received err in TestNewPGListener: %v", ERR)
	}
}

func callBackTest(p *pq.Notification) {
	CALLBACKGLOBAL = p.Extra
}

func TestBeginAndEndListen(t *testing.T) {
	go LISTENER.BeginListen(nil)
	LISTENER.killListen <- false
}

func TestNotifyUpdate(t *testing.T) {
	go LISTENER.BeginListen(callBackTest)

	dbConn, err := OpenConnectionWithConfig(&CONFIG)
	if err != nil {
		t.Fatalf("Error encountered in TestNotifyUpdate: %v", err)
	}

	whiteTorrent := schemas.WhiteTorrent{
		InfoHash:  "InfoHash",
		Name:      "Test NOTIFY",
		AddedBy:   "alksdjfj",
		DateAdded: 0,
	}

	retval := whiteTorrent.AddWhitelistedTorrent(dbConn)
	if retval == false {
		t.Fatalf("Failed to insert white torrent in TestNotifyUpdate")
	}

	time.Sleep(5 * time.Millisecond)

	if CALLBACKGLOBAL == "" {
		t.Fatalf("Got %s", CALLBACKGLOBAL)
	}
}

func TestMain(m *testing.M) {
	dbConn, _ := OpenConnectionWithConfig(&CONFIG)
	dbConn.DropTableIfExists(
		&schemas.PeerStats{},
		&schemas.Torrent{},
		&schemas.TrackerStats{},
	)
	InitDB(dbConn)
	os.Exit(m.Run())

}

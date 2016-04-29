package server

import (
	"github.com/GrappigPanda/notorious/database"
	"testing"
)

var DBCONN, _ = db.OpenConnection()

var CONTEXT = requestAppContext{
	dbConn:      DBCONN,
	redisClient: OpenClient(),
	whitelist:   true,
}

var ANNOUNCEDATA = &announceData{
	info_hash:      "12345123451234512345",
	peer_id:        "11111111111111111111",
	ip:             "127.0.0.1",
	event:          "STARTED",
	port:           6667,
	uploaded:       0,
	downloaded:     0,
	left:           0,
	numwant:        30,
	compact:        true,
	requestContext: CONTEXT,
}

// TestStartedEventHandler tests that with a whitelist being active, we can not
// add a new info_hash to the tracker.
func TestStartedEventHandler(t *testing.T) {
	err := ANNOUNCEDATA.StartedEventHandler()

	if err == nil {
		t.Fatalf("Failed to TestStartedEventHandler: %v", err)
	}
}

// TestStartedEventHandler tests that with a whitelist being active, we can not
// add a new info_hash to the tracker.
func TestStartedEventHandlerNoWhitelist(t *testing.T) {
	announce2 := ANNOUNCEDATA
	announce2.requestContext.whitelist = false
	err := announce2.StartedEventHandler()

	if err != nil {
		t.Fatalf("Failed to TestStartedEventHandler: %v", err)
	}
}

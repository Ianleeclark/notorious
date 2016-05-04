package db

import (
	"testing"
	"time"
	"fmt"
)

var DBCONN, _ = OpenConnection()

func TestOpenConn(t *testing.T) {
	dbConn, err := OpenConnection()
	if err != nil {
		t.Fatalf("%v", err)
	}
	InitDB(dbConn)
}

func TestAddWhitelistedTorrent(t *testing.T) {
	newTorrent := &White_Torrent{
		InfoHash:   "12345123451234512345",
		Name:       "Hello Kitty Island Adventure.exe",
		AddedBy:    "127.0.0.1",
		DateAdded:  time.Now().Unix(),
	}

	if !newTorrent.AddWhitelistedTorrent() {
		t.Fatalf("Failed to Add a whitelisted torrent")
	}
}

func TestGetWhitelistedTorrents(t *testing.T) {
	newTorrent := &White_Torrent{
		InfoHash:   "12345123GetWhitelistedTorrents",
		Name:       "Hello Kitty Island Adventure3.exe",
		AddedBy:    "127.0.0.1",
		DateAdded:  time.Now().Unix(),
	}

	newTorrent.AddWhitelistedTorrent()

	retval, err := GetWhitelistedTorrents()
	if err != nil {
		t.Fatalf("Failed to get all whitelisted torrents: %v", err)
	}

	fmt.Printf("VALUES: %v", retval)
}

func TestGetWhitelistedTorrent(t *testing.T) {
	newTorrent := &White_Torrent{
		InfoHash:   "12345123GetWhitelistedTorrent",
		Name:       "Hello Kitty Island Adventure2.exe",
		AddedBy:    "127.0.0.1",
		DateAdded:  time.Now().Unix(),
	}

	newTorrent.AddWhitelistedTorrent()

	retval, err := GetWhitelistedTorrent(newTorrent.InfoHash)
	if err != nil {
		t.Fatalf("Failed to GetWhitelistedTorrent: %v", err)
	}

	if retval.InfoHash != newTorrent.InfoHash {
		t.Fatalf("Expected %v, got %v", retval.InfoHash,
			newTorrent.InfoHash)
	}
}

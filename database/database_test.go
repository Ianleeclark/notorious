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
	newTorrent := &Torrent{
		InfoHash:   "12345123451234512345",
		Name:       "Hello Kitty Island Adventure.exe",
		Downloaded: 0,
		Seeders:    0,
		Leechers:   0,
		AddedBy:    "127.0.0.1",
		DateAdded:  time.Now().Unix(),
	}

	newTorrent.AddWhitelistedTorrent()

	retval, err := GetTorrent(newTorrent.InfoHash)
	if err != nil {
		t.Fatalf("Failed to GetTorrent")
	}

	if newTorrent.InfoHash != retval.InfoHash {
		t.Fatalf("Expected %v, got %v", newTorrent.InfoHash, retval.InfoHash)
	}

	if newTorrent.DateAdded != retval.DateAdded {
		t.Fatalf("Expected %v, got %v", newTorrent.DateAdded, retval.DateAdded)
	}

	if newTorrent.Name != retval.Name {
		t.Fatalf("Expected %v, got %v", newTorrent.Name, retval.Name)
	}
}

func TestGetWhitelistedTorrents(t *testing.T) {
	newTorrent := &Torrent{
		InfoHash:   "12345123GetWhitelistedTorrents",
		Name:       "Hello Kitty Island Adventure3.exe",
		Downloaded: 0,
		Seeders:    0,
		Leechers:   0,
		AddedBy:    "127.0.0.1",
		DateAdded:  time.Now().Unix(),
	}

	newTorrent.AddWhitelistedTorrent()
	// Done just to verify the error isn't with getting whitelisted torrents
	_, err := GetTorrent(newTorrent.InfoHash)
 	if err != nil {
 		t.Fatalf("Failed to GetTorrent")
 	}
	
	retval, err := GetWhitelistedTorrents()
	if err != nil {
		t.Fatalf("Failed to get all whitelisted torrents: %v", err)
	}

	fmt.Printf("%v", retval)
}

func TestGetWhitelistedTorrent(t *testing.T) {
	newTorrent := &Torrent{
		InfoHash:   "12345123GetWhitelistedTorrent",
		Name:       "Hello Kitty Island Adventure2.exe",
		Downloaded: 0,
		Seeders:    0,
		Leechers:   0,
		AddedBy:    "127.0.0.1",
		DateAdded:  time.Now().Unix(),
	}

	newTorrent.AddWhitelistedTorrent()
	// Done just to verify the error isn't with getting whitelisted torrents
	_, err := GetTorrent(newTorrent.InfoHash)
 	if err != nil {
 		t.Fatalf("Failed to GetTorrent")
 	}
 	
	retval, err := GetWhitelistedTorrent(newTorrent.InfoHash)
	if err != nil {
		t.Fatalf("Failed to GetWhitelistedTorrent: %v", err)
	}

	if retval.InfoHash != newTorrent.InfoHash {
		t.Fatalf("Expected %v, got %v", retval.InfoHash,
			newTorrent.InfoHash)
	}
}

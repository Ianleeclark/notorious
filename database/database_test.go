package db

import (
    "time"
    "testing"
)

var DBCONN, _ = OpenConnection()

func TestOpenConn(t *testing.T) {
    _, err := OpenConnection()
    if err != nil {
        t.Fatalf("%v", err)
    }
}

func TestAddWhitelistedTorrent(t *testing.T) {
    newTorrent := Torrent{
        infoHash: "12345123451234512345",
        name: "Hello Kitty Island Adventure.exe",
        downloaded: 0,
        seeders: 0,
        leechers: 0,
        addedBy: "127.0.0.1",
        dateAdded: time.Now().UTC(),
    }

    newTorrent.AddWhitelistedTorrent()

    retval, err := GetTorrent(newTorrent.infoHash)
    if err != nil {
        t.Fatalf("Failed to GetTorrent %v", err)
    }

    if retval.infoHash != newTorrent.infoHash {
        t.Fatalf("Expected %v, got %v", retval.infoHash, newTorrent.infoHash)
    }
}

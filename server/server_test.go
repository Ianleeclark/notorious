package server

import (
	"testing"
)

func TestParseTorrentGetRequestURI(t *testing.T) {
	expectedResult := map[string]interface{}{
		"info_hash": "1",
		"peer_id":   "5",
		"ip":        "127.0.0.1",
	}
	result := parseTorrentGetRequestURI("/announce?info_hash=1%26peer_id=5%26ip=127.0.0.1")
	if result["info_hash"] != expectedResult["info_hash"] {
		t.Fatalf("Expected %s got %s", expectedResult, result)
	}
	if result["peer_id"] != expectedResult["peer_id"] {
		t.Fatalf("Expected %s got %s", expectedResult, result)
	}
	if result["ip"] != expectedResult["ip"] {
		t.Fatalf("Expected %s got %s", expectedResult, result)
	}
}

func TestFillEmptyMapValues(t *testing.T) {
	expectedResult := TorrentRequestData{"1", "5", "127.0.0.1", 0, 0, 0, 0, STOPPED}

	x := parseTorrentGetRequestURI("/announce?info_hash=1%26peer_id=5%26ip=127.0.0.1")

	result := fillEmptyMapValues(x)
	if result.port != 0 {
		t.Fatalf("Expected %s got %s", expectedResult.port, result.port)
	}
	if result.uploaded != 0 {
		t.Fatalf("Expected %s got %s", expectedResult.uploaded, result.uploaded)
	}
	if result.downloaded != 0 {
		t.Fatalf("Expected %s got %s", expectedResult.downloaded, result.downloaded)
	}
	if result.left != 0 {
		t.Fatalf("Expected %s got %s", expectedResult.left, result.left)
	}
	if result.event != STOPPED {
		t.Fatalf("Expected %s got %s", expectedResult.event, result.event)
	}

}

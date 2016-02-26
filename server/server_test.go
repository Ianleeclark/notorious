package server

import (
	"fmt"
	"testing"
)

func TestParseUrlQuery(t *testing.T) {
	request := "http://127.0.0.1:3000/announce?info_hash=QtA%C0%81%8D%C5GV%02%150%5D%2B%91%80a%BB%02%9A&peer_id=-lt0D20-s%081%8ER%D7%C9%15X%DB%DD%D2&key=602bcd6f&compact=1&port=6963&uploaded=0&downloaded=0&left=5448254&event=started"

	result := decodeQueryURL(request)
	if result["uploaded"][0] != "0" {
		t.Fatalf("Expected 0, got %s", result["uploaded"])
	}
	if result["port"][0] != "6963" {
		t.Fatalf("Expected 0, got %s", result["port"])
	}
	if result["downloaded"][0] != "0" {
		t.Fatalf("Expected 0, got %s", result["downloaded"])
	}
	if result["compact"][0] != "1" {
		t.Fatalf("Expected 0, got %s", result["compact"])
	}

}

func TestParseTorrentGetRequest(t *testing.T) {
	request := "http://127.0.0.1:3000/announce?info_hash=QtA%C0%81%8D%C5GV%02%150%5D%2B%91%80a%BB%02%9A&peer_id=-lt0D20-s%081%8ER%D7%C9%15X%DB%DD%D2&key=602bcd6f&compact=1&port=6963&uploaded=0&downloaded=0&left=5448254&event=started"

	result := decodeQueryURL(request)
	fmt.Println(result)
	if result["uploaded"][0] != "0" {
		t.Fatalf("Expected 0, got %s", result["uploaded"])
	}
	if result["port"][0] != "6963" {
		t.Fatalf("Expected 0, got %s", result["port"])
	}
	if result["downloaded"][0] != "0" {
		t.Fatalf("Expected 0, got %s", result["downloaded"])
	}
	if result["compact"][0] != "1" {
		t.Fatalf("Expected 0, got %s", result["compact"])
	}
}

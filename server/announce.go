package server

import (
	"fmt"
	"net/url"
)

type announceData struct {
	info_hash  string //20 byte sha1 hash
	peer_id    string //max len 20
	ip         string //optional
	port       string // port number the peer is listening on
	uploaded   int    // base10 ascii amount uploaded so far
	downloaded int    // base10 ascii amount downloaded so far
	left       int    // # of bytes left to download (base 10 ascii)
	event      int    // TorrentEvent
	numwant    int    // Number of peers requested by client.
	compact    bool   // Bep23 peer list compression decision: True -> compress bep23
}

func (a *announceData) parseAnnounceData(u *url.URL) {

}

func ParseInfoHash(s string) string {
	return fmt.Sprintf("%x", s)
}

func decodeQueryURL(s string) url.Values {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	m, _ := url.ParseQuery(u.RawQuery)
	return m
}

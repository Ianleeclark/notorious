package server

import (
	"fmt"
	"net/url"
	"strconv"
)

type announceData struct {
	info_hash  string //20 byte sha1 hash
	peer_id    string //max len 20
	ip         string //optional
	event      string // TorrentEvent
	port       uint64 // port number the peer is listening on
	uploaded   uint64 // base10 ascii amount uploaded so far
	downloaded uint64 // base10 ascii amount downloaded so far
	left       uint64 // # of bytes left to download (base 10 ascii)
	numwant    uint64 // Number of peers requested by client.
	compact    bool   // Bep23 peer list compression decision: True -> compress bep23
}

func (a *announceData) parseAnnounceData(u *url.URL) (err error) {
	query := u.Query()
	a.info_hash = ParseInfoHash(query.Get("info_hash"))
	if a.info_hash == "" {
		err = fmt.Errorf("No info_hash provided.")
		return
	}
	a.ip = query.Get("ip")
	if a.ip == "" {
		return fmt.Errorf("No info_hash provided.")
	}
	a.peer_id = query.Get("peer_id")
	if a.peer_id == "" {
		return fmt.Errorf("No info_hash provided.")
	}
	a.port, err = GetInt(query, "port")
	if err != nil {
		return fmt.Errorf("Failed to get port")
	}
	a.downloaded, err = GetInt(query, "downloaded")
	if err != nil {
		err = fmt.Errorf("Failed to get downloaded byte count.")
		return
	}
	a.uploaded, err = GetInt(query, "uploaded")
	if err != nil {
		err = fmt.Errorf("Failed to get uploaded byte count.")
		return
	}
	a.left, err = GetInt(query, "left")
	if err != nil {
		err = fmt.Errorf("Failed to get remaining byte count.")
		return
	}
	a.numwant, err = GetInt(query, "numwant")
	if err != nil {
		err = fmt.Errorf("Failed to get number of peers requested.")
		return
	}
	if x := query.Get("compact"); x != "" {
		a.compact, err = strconv.ParseBool(x)
		if err != nil {
			err = fmt.Errorf("Failed to parse a boolean value from `compact`.")
			return
		}
	}
	a.event = query.Get("event")

	return
}

func GetInt(u url.Values, key string) (ui uint64, err error) {
	if x := u.Get(key); x == "" {
		err = fmt.Errorf("Failed to locate the key in the url.")
	} else {
		ui, err = strconv.ParseUint(x, 10, 64)
		if err != nil {
			err = fmt.Errorf("Failed to parse uint from the key")
			return
		}
		return
	}
	return
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

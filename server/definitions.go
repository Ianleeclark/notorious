package server

import (
	"net"
	"time"
)

const (
	STARTED = iota
	COMPLETED
	STOPPED
)

type TorrentEvent struct {
	started   int
	completed int
	stopped   int
}

type TorrentRequestData struct {
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

type Peer struct {
	remoteAddr  *net.TCPAddr // Remote connection address
	downloaded  uint64       // Total bytes downloaded
	uploaded    uint64       // Total bytes uploaded
	left        uint64       // Bytes left for torrent
	lastTracked time.Time    // Last known appearance in the tracker
}

type PeerList map[string]*Peer

type TorrentResponseData struct {
	interval     int
	min_interval int
	tracker_id   string
	completed    int
	incomplete   int
	peers        interface{}
}

var ANNOUNCE_URL = "/announce"

package server

import (
	"github.com/NotoriousTracker/redis"
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

type announceData struct {
	info_hash   string        //20 byte sha1 hash
	peer_id     string        //max len 20
	ip          string        //optional
	event       string        // TorrentEvent
	port        uint64        // port number the peer is listening on
	uploaded    uint64        // base10 ascii amount uploaded so far
	downloaded  uint64        // base10 ascii amount downloaded so far
	left        uint64        // # of bytes left to download (base 10 ascii)
	numwant     uint64        // Number of peers requested by client.
	compact     bool          // Bep23 peer list compression decision: True -> compress bep23
	redisClient *redis.Client // The redis client connection handler to use.
}

type TorrentResponseData struct {
	interval     int
	min_interval int
	tracker_id   string
	completed    int
	incomplete   int
	peers        interface{}
}

var ANNOUNCE_URL = "/announce"

// TODO(ian): Set this expireTime to a config-loaded value.
// expireTime := 5 * 60

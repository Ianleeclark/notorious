package server

import (
	"github.com/GrappigPanda/notorious/config"
	"github.com/jinzhu/gorm"
	"gopkg.in/redis.v3"
)

const (
	RATIOLESS = iota
	SEMIRATIOLESS
	NORMALRATIO
)

type announceData struct {
	info_hash      string            //20 byte sha1 hash
	peer_id        string            //max len 20
	ip             string            //optional
	event          string            // TorrentEvent
	port           uint64            // port number the peer is listening on
	uploaded       uint64            // base10 ascii amount uploaded so far
	downloaded     uint64            // base10 ascii amount downloaded so far
	left           uint64            // # of bytes left to download (base 10 ascii)
	numwant        uint64            // Number of peers requested by client.
	compact        bool              // Bep23 peer list compression decision: True -> compress bep23
	requestContext requestAppContext // The request-specific connections
}

// requestAppContext First of all naming things is the hardest part of
// programming real talk. Second of all, this essentially houses
// request-specific data like db connections and in the future the redisClient.
// Things that should persist only within the duration of a request.
type requestAppContext struct {
	dbConn      *gorm.DB
	redisClient *redis.Client // The redis client connection handler to use.
	whitelist   bool
}

type scrapeData struct {
	infoHash string
}

// scrapeResponse is the data associated with a returned scrape.
type scrapeResponse struct {
	complete   uint64
	downloaded uint64
	incomplete uint64
}

// TorrentResponseData models what is sent back to the peer upon a succesful
// info hash lookup.
type TorrentResponseData struct {
	interval     int
	min_interval int
	tracker_id   string
	completed    int
	incomplete   int
	peers        interface{}
}

// ANNOUNCE_URL The announce path for the http calls to reach.
var ANNOUNCE_URL = "/announce"

// TODO(ian): Set this expireTime to a config-loaded value.
// expireTime := 5 * 60

// applicationContext houses data necessary for the handler to properly
// function as the application is desired.
type applicationContext struct {
	config       config.ConfigStruct
	trackerLevel int
}

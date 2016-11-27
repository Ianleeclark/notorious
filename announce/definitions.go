package announce

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/redis.v3"
)

const (
	RATIOLESS = iota
	SEMIRATIOLESS
	NORMALRATIO
)

type AnnounceData struct {
	InfoHash string //20 byte sha1 hash
	PeerID   string //max len 20
	IP       string //optional
	Event    string // TorrentEvent

	Port uint64 // port number the peer is listening
	// on

	Uploaded   uint64 // base10 ascii amount uploaded so far
	Downloaded uint64 // base10 ascii amount downloaded so
	// far

	Left uint64 // # of bytes left to download
	// (base 10 ascii)

	Numwant uint64 // Number of peers requested by client.

	Compact bool // Bep23 peer list compression
	// decision: True -> compress bep23

	RequestContext requestAppContext // The request-specific connections
}

// requestAppContext First of all naming things is the hardest part of
// programming real talk. Second of all, this essentially houses
// request-specific data like db connections and in the future the redisClient.
// Things that should persist only within the duration of a request.
type requestAppContext struct {
	dbConn      *gorm.DB
	redisClient *redis.Client // The redis client connection handler to use.
	Whitelist   bool
}

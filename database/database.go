package db

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	// We use a blank import here because I'm afraid of breaking anything
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// AddWhitelistedTorrent handles adding a new whitelisted torrent into the
// database.
func (t *WhiteTorrent) AddWhitelistedTorrent(db *gorm.DB) bool {
	db.Create(t)
	return db.NewRecord(t)
}

// PeerDeltaEvent allows us to set an event to handle how we're going to update
// teh database.
type PeerDeltaEvent int

const (
	// PEERUPDATE represents a change to a peer, so we'll update a tracker
	// user's ratio.
	PEERUPDATE PeerDeltaEvent = iota
	// TRACKERUPDATE handles updating total tracker stats
	TRACKERUPDATE
	// TORRENTUPDATE represents the changes to a specific torrent where we
	// update total upload/download for the torrent itself.
	TORRENTUPDATE
)

// PeerTrackerDelta handles holding data to be updated by the `UpdateConsumer`.
type PeerTrackerDelta struct {
	Uploaded   uint64
	Downloaded uint64
	IP         string
	Event      PeerDeltaEvent
}

// SQLStore is the base implementation for a database which will be used to
// store stats and retrieve whitelisted torrents.
type SQLStore interface {
	OpenConnection() (*gorm.DB, error)
	GetTorrent(string) (*Torrent, error)
	GetWhitelistedTorrent(string) (*WhiteTorrent, error)
	UpdateStats(uint64, uint64)
	UpdateTorrentStats(int64, int64)
	ScrapeTorrent(string) *Torrent
	GetWhitelistedTorrents() (*sql.Rows, error)
	UpdatePeerStats(uint64, uint64, string)
	HandlePeerUpdates() chan PeerTrackerDelta
}

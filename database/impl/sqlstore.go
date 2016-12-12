package sqlStoreImpl

import (
	"database/sql"
	"github.com/GrappigPanda/notorious/config"
	"github.com/GrappigPanda/notorious/database/schemas"
	"github.com/jinzhu/gorm"
)

// SQLStore is the base implementation for a database which will be used to
// store stats and retrieve whitelisted torrents.
type SQLStore interface {
	OpenConnection() (*gorm.DB, error)
	GetTorrent(string) (*schemas.Torrent, error)
	GetWhitelistedTorrent(string) (*schemas.WhiteTorrent, error)
	UpdateStats(uint64, uint64)
	UpdateTorrentStats(int64, int64)
	ScrapeTorrent(string) *schemas.Torrent
	GetWhitelistedTorrents() (*sql.Rows, error)
	UpdatePeerStats(uint64, uint64, string)
	HandlePeerUpdates() chan PeerTrackerDelta
}

func InitSQLStoreByDBChoice() SQLStore {
	cfg := config.LoadConfig()
	if cfg.DBChoice == "mysql" {
		return new(MySQLStore)
	} else if cfg.DBChoice == "postgres" {
		return new(PostgresStore)
	} else {
		panic("Invalid database choice found for `InitSQLStoreByDBChoice`.")
	}
}

// PeerTrackerDelta handles holding data to be updated by the `UpdateConsumer`.
type PeerTrackerDelta struct {
	Uploaded   uint64
	Downloaded uint64
	IP         string
	Event      PeerDeltaEvent
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

package db

import (
	"database/sql"
	"github.com/GrappigPanda/notorious/config"
	"github.com/GrappigPanda/notorious/database/mysql"
	"github.com/GrappigPanda/notorious/database/postgres"
	"github.com/GrappigPanda/notorious/database/schemas"
	"github.com/jinzhu/gorm"
	// We use a blank import here because I'm afraid of breaking anything
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func InitDB(config *config.ConfigStruct) {
	if (*config).DBChoice == "mysql" {
		conn, err := mysql.OpenConnection()
		if err != nil {
			panic("Unable to open connection to remote server")
		}
		mysql.InitDB(conn)
	} else if (*config).DBChoice == "postgres" {
		conn, err := postgres.OpenConnection()
		if err != nil {
			panic("Unable to open connection to remote server")
		}
		postgres.InitDB(conn)
	} else {
		panic("Invalid Config choice for DBChoice. Set either `UsePostgres` or `UseMySQL`.")
	}
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
	OpenConnectionWithConfig(*config.ConfigStruct) (*gorm.DB, error)
	GetTorrent(string) (*schemas.Torrent, error)
	GetWhitelistedTorrent(string) (*schemas.WhiteTorrent, error)
	UpdateStats(uint64, uint64)
	UpdateTorrentStats(int64, int64)
	ScrapeTorrent(string) *schemas.Torrent
	GetWhitelistedTorrents() (*sql.Rows, error)
	UpdatePeerStats(uint64, uint64, string)
	HandlePeerUpdates() chan PeerTrackerDelta
}

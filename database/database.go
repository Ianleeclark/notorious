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
}

package sqlStoreImpl

import (
	"database/sql"
	"github.com/GrappigPanda/notorious/database"
	"github.com/GrappigPanda/notorious/database/mysql"
	"github.com/jinzhu/gorm"
	// We use a blank import here because I'm afraid of breaking anything
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// MySQLStore represents the mysql implementation of `SQLStore`
type MySQLStore struct {
	dbPool *gorm.DB
}

// OpenConnection wraps `mysql.OpenConnection`.
func (m *MySQLStore) OpenConnection() (*gorm.DB, error) {
	return mysql.OpenConnection()
}

// GetTorrent wraps `mysql.GetTorrent`.
func (m *MySQLStore) GetTorrent(infoHash string) (*db.Torrent, error) {
	return mysql.GetTorrent(m.dbPool, infoHash)
}

// GetWhitelistedTorrent wraps `mysql.GetWhitelistedTorrent`.
func (m *MySQLStore) GetWhitelistedTorrent(infoHash string) (*db.White_Torrent, error) {
	return mysql.GetWhitelistedTorrent(m.dbPool, infoHash)
}

// UpdateStats wraps `mysql.UpdateStats`.
func (m *MySQLStore) UpdateStats(uploaded uint64, downloaded uint64) {
	mysql.UpdateStats(m.dbPool, uploaded, downloaded)
}

// UpdateTorrentStats wraps `mysql.UpdateTorrentStats`.
func (m *MySQLStore) UpdateTorrentStats(uploaded int64, downloaded int64) {
	mysql.UpdateTorrentStats(m.dbPool, uploaded, downloaded)
}

// ScrapeTorrent wraps `mysql.ScrapeTorrent`.
func (m *MySQLStore) ScrapeTorrent(infoHash string) *db.Torrent {
	return mysql.ScrapeTorrent(m.dbPool, infoHash)
}

// GetWhitelistedTorrents wraps `mysql.GetWhitelistedTorrents`.
func (m *MySQLStore) GetWhitelistedTorrents() (*sql.Rows, error) {
	return mysql.GetWhitelistedTorrents(m.dbPool)
}

// UpdatePeerStats wraps `mysql.UpdatePeerStats`.
func (m *MySQLStore) UpdatePeerStats(uint64, uint64, string) {

}

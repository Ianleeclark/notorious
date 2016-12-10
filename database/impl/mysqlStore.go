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
	dbPool         *gorm.DB
	UpdateConsumer chan db.PeerTrackerDelta
}

// InitMySQLStore Creates a `MySQLStore` object and initiates all necessary
// moving parts like the `HandlePeerUpdates` channel consumer.
func InitMySQLStore() (store MySQLStore) {
	dbConn, err := mysql.OpenConnection()
	if err != nil {
		panic("Failed opening a connection to remote MYSQL database.")
	}

	store = MySQLStore{
		dbConn,
		nil,
	}

	store.UpdateConsumer = store.HandlePeerUpdates()

	return store
}

// OpenConnection wraps `mysql.OpenConnection`.
func (m *MySQLStore) OpenConnection() (*gorm.DB, error) {
	return mysql.OpenConnection()
}

// HandlePeerUpdates handles listening and aggregating peer updates. THis
// allows block/asynchronous consumption of peer updates, rather than updating
// the remote database at the end of every request.
func (m *MySQLStore) HandlePeerUpdates() chan db.PeerTrackerDelta {
	peerUpdatesChan := make(chan db.PeerTrackerDelta)

	go func() {
		for {
			update := <-peerUpdatesChan
			switch update.Event {
			case db.PEERUPDATE:
				m.UpdatePeerStats(update.Uploaded, update.Downloaded, update.IP)
			case db.TRACKERUPDATE:
				m.UpdateStats(update.Uploaded, update.Downloaded)
			case db.TORRENTUPDATE:
				m.UpdateTorrentStats(int64(update.Uploaded), int64(update.Downloaded))
			}
		}
	}()

	return peerUpdatesChan
}

// GetTorrent wraps `mysql.GetTorrent`.
func (m *MySQLStore) GetTorrent(infoHash string) (*db.Torrent, error) {
	return mysql.GetTorrent(m.dbPool, infoHash)
}

// GetWhitelistedTorrent wraps `mysql.GetWhitelistedTorrent`.
func (m *MySQLStore) GetWhitelistedTorrent(infoHash string) (*db.WhiteTorrent, error) {
	return mysql.GetWhitelistedTorrent(m.dbPool, infoHash)
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
func (m *MySQLStore) UpdatePeerStats(uploaded uint64, downloaded uint64, ip string) {
	mysql.UpdatePeerStats(m.dbPool, uploaded, downloaded, ip)
}

// UpdateStats wraps `mysql.UpdateStats`.
func (m *MySQLStore) UpdateStats(uploaded uint64, downloaded uint64) {
	m.UpdateConsumer <- db.PeerTrackerDelta{
		Uploaded:   uploaded,
		Downloaded: downloaded,
	}
	mysql.UpdateStats(m.dbPool, uploaded, downloaded)
}

// UpdateTorrentStats wraps `mysql.UpdateTorrentStats`.
func (m *MySQLStore) UpdateTorrentStats(uploaded int64, downloaded int64) {
	mysql.UpdateTorrentStats(m.dbPool, uploaded, downloaded)
}

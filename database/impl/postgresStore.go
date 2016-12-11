package sqlStoreImpl

import (
	"database/sql"
	"github.com/GrappigPanda/notorious/database"
	"github.com/GrappigPanda/notorious/database/mysql"
	"github.com/GrappigPanda/notorious/database/schemas"
	"github.com/jinzhu/gorm"
	// We use a blank import here because I'm afraid of breaking anything
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type PostgresStore struct {
	dbPool *gorm.DB
	// NOTE: The `UpdateConsumer` as seen in `mysqlStore.go:PostgresStore is
	// unnecessary in Postgres, as we will rely on the `pg_notify` feature
	// implemented in postgres.
}

// InitPostgresStore Creates a `PostgresStore` object and initiates all necessary
// moving parts like the `HandlePeerUpdates` channel consumer.
func InitPostgresStore() (store PostgresStore) {
	dbConn, err := mysql.OpenConnection()
	if err != nil {
		panic("Failed opening a connection to remote Postgres database.")
	}

	store = PostgresStore{
		dbConn,
	}

	return store
}

// OpenConnection wraps `mysql.OpenConnection`.
func (m *PostgresStore) OpenConnection() (*gorm.DB, error) {
	return mysql.OpenConnection()
}

// HandlePeerUpdates handles listening and aggregating peer updates. THis
// allows block/asynchronous consumption of peer updates, rather than updating
// the remote database at the end of every request.
func (m *PostgresStore) HandlePeerUpdates() chan db.PeerTrackerDelta {
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
func (m *PostgresStore) GetTorrent(infoHash string) (*schemas.Torrent, error) {
	return mysql.GetTorrent(m.dbPool, infoHash)
}

// GetWhitelistedTorrent wraps `mysql.GetWhitelistedTorrent`.
func (m *PostgresStore) GetWhitelistedTorrent(infoHash string) (*schemas.WhiteTorrent, error) {
	return mysql.GetWhitelistedTorrent(m.dbPool, infoHash)
}

// ScrapeTorrent wraps `mysql.ScrapeTorrent`.
func (m *PostgresStore) ScrapeTorrent(infoHash string) *schemas.Torrent {
	return mysql.ScrapeTorrent(m.dbPool, infoHash)
}

// GetWhitelistedTorrents wraps `mysql.GetWhitelistedTorrents`.
func (m *PostgresStore) GetWhitelistedTorrents() (*sql.Rows, error) {
	return mysql.GetWhitelistedTorrents(m.dbPool)
}

// UpdatePeerStats wraps `mysql.UpdatePeerStats`.
func (m *PostgresStore) UpdatePeerStats(uploaded uint64, downloaded uint64, ip string) {
	mysql.UpdatePeerStats(m.dbPool, uploaded, downloaded, ip)
}

// UpdateStats wraps `mysql.UpdateStats`.
func (m *PostgresStore) UpdateStats(uploaded uint64, downloaded uint64) {
	mysql.UpdateStats(m.dbPool, uploaded, downloaded)
}

// UpdateTorrentStats wraps `mysql.UpdateTorrentStats`.
func (m *PostgresStore) UpdateTorrentStats(uploaded int64, downloaded int64) {
	mysql.UpdateTorrentStats(m.dbPool, uploaded, downloaded)
}

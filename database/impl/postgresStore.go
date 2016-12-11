package sqlStoreImpl

import (
	"database/sql"
	"github.com/GrappigPanda/notorious/database"
	"github.com/GrappigPanda/notorious/database/postgres"
	"github.com/GrappigPanda/notorious/database/schemas"
	"github.com/jinzhu/gorm"
	//"github.com/lib/pq"
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
	dbConn, err := postgres.OpenConnection()
	if err != nil {
		panic("Failed opening a connection to remote Postgres database.")
	}

	store = PostgresStore{
		dbConn,
	}

	return store
}

// OpenConnection wraps `postgres.OpenConnection`.
func (m *PostgresStore) OpenConnection() (*gorm.DB, error) {
	return postgres.OpenConnection()
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

// GetTorrent wraps `postgres.GetTorrent`.
func (m *PostgresStore) GetTorrent(infoHash string) (*schemas.Torrent, error) {
	return postgres.GetTorrent(m.dbPool, infoHash)
}

// GetWhitelistedTorrent wraps `postgres.GetWhitelistedTorrent`.
func (m *PostgresStore) GetWhitelistedTorrent(infoHash string) (*schemas.WhiteTorrent, error) {
	return postgres.GetWhitelistedTorrent(m.dbPool, infoHash)
}

// ScrapeTorrent wraps `postgres.ScrapeTorrent`.
func (m *PostgresStore) ScrapeTorrent(infoHash string) *schemas.Torrent {
	return postgres.ScrapeTorrent(m.dbPool, infoHash)
}

// GetWhitelistedTorrents wraps `postgres.GetWhitelistedTorrents`.
func (m *PostgresStore) GetWhitelistedTorrents() (*sql.Rows, error) {
	return postgres.GetWhitelistedTorrents(m.dbPool)
}

// UpdatePeerStats wraps `postgres.UpdatePeerStats`.
func (m *PostgresStore) UpdatePeerStats(uploaded uint64, downloaded uint64, ip string) {
	postgres.UpdatePeerStats(m.dbPool, uploaded, downloaded, ip)
}

// UpdateStats wraps `postgres.UpdateStats`.
func (m *PostgresStore) UpdateStats(uploaded uint64, downloaded uint64) {
	postgres.UpdateStats(m.dbPool, uploaded, downloaded)
}

// UpdateTorrentStats wraps `postgres.UpdateTorrentStats`.
func (m *PostgresStore) UpdateTorrentStats(uploaded int64, downloaded int64) {
	postgres.UpdateTorrentStats(m.dbPool, uploaded, downloaded)
}

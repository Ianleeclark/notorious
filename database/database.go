package db

import (
	"database/sql"
	"fmt"
	"github.com/GrappigPanda/notorious/config"
	"github.com/jinzhu/gorm"
	// We use a blank import here because I'm afraid of breaking anything
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// OpenConnection does as its name dictates and opens a connection to the
// MysqlHost listed in the config
func OpenConnection() (db *gorm.DB, err error) {
	c := config.LoadConfig()

	db, err = gorm.Open("mysql", formatConnectString(c))
	if err != nil {
		err = fmt.Errorf("Failed to open connection to MySQL: %v", err)
	}

	return
}

// InitDB initializes database tables.
func InitDB(db *gorm.DB) {
	db = assertOpenConnection(db)

	db.CreateTable(&White_Torrent{})
	db.CreateTable(&Torrent{})
	db.CreateTable(&TrackerStats{})
	db.CreateTable(&Peer_Stats{})
}

// AddWhitelistedTorrent adds a torrent to the whitelist so that they may be
// used by the tracker in the future.
func (t *White_Torrent) AddWhitelistedTorrent(db *gorm.DB) bool {
	db = assertOpenConnection(db)

	db.Create(t)
	return db.NewRecord(t)
}

// GetTorrent retrieves a torrent by its infoHash from the generic torrent
// table in the database. Note: there's also a whitelisted torrent table
// (`white_torrent`).
func GetTorrent(infoHash string) (db *gorm.DB, t *Torrent, err error) {
	db = assertOpenConnection(db)

	t = &Torrent{}

	db.Where("info_hash = ?", infoHash).Find(&t)

	return
}

// GetWhitelistedTorrent Retrieves a single whitelisted torrent by its infoHash
func GetWhitelistedTorrent(infoHash string) (db *gorm.DB, t *White_Torrent, err error) {
	db = assertOpenConnection(db)

	t = &White_Torrent{}

	x := db.Where("info_hash = ?", infoHash).First(&t)
	if x.Error != nil {
		err = x.Error
	}

	return
}

// UpdateStats Handles updating statistics relevant to our tracker.
func UpdateStats(db *gorm.DB, uploaded uint64, downloaded uint64) {
	db = assertOpenConnection(db)

	ts := &TrackerStats{}
	db.First(&ts)
	db.Model(&ts).Updates(
		TrackerStats{
			Uploaded:   ts.Uploaded + int64(uploaded),
			Downloaded: ts.Downloaded + int64(downloaded),
		})

	return
}

// UpdateTorrentStats Handles updating statistics relevant to our tracker.
func UpdateTorrentStats(db *gorm.DB, seederDelta int64, leecherDelta int64) {
	db = assertOpenConnection(db)

	t := &Torrent{}
	db.First(&t)
	db.Model(&t).Updates(
		TrackerStats{
			Uploaded:   t.Seeders + seederDelta,
			Downloaded: t.Leechers + leecherDelta,
		})

	return
}

// UpdatePeerStats handles updating peer info like hits per ip, downloaded
// amount, uploaded amounts.
func UpdatePeerStats(db *gorm.DB, uploaded uint64, downloaded uint64, ip string) {
	db = assertOpenConnection(db)

	ps := &Peer_Stats{Ip: ip}
	db.First(&ps)
	db.Model(&ps).UpdateColumn(map[string]interface{}{
		"Uploaded":   ps.Uploaded + int64(uploaded),
		"Downloaded": ps.Downloaded + int64(downloaded),
	})

	return
}

// GetWhitelistedTorrents allows us to retrieve all of the white listed
// torrents. Mostly used for populating the Redis KV storage with all of our
// whitelisted torrents.
func GetWhitelistedTorrents(db *gorm.DB) (x *sql.Rows, err error) {
	db = assertOpenConnection(db)

	x, err = db.Table("white_torrents").Rows()
	if err != nil {
		return
	}

	return
}

// ScrapeTorrent supports the Scrape convention
func ScrapeTorrent(db *gorm.DB, infoHash string) (torrent *Torrent) {
	db = assertOpenConnection(db)

	db.Where("info_hash = ?", infoHash).First(&torrent)
	return
}

// formatConnectStrings concatenates the data from the config file into a
// usable MySQL connection string.
func formatConnectString(c config.ConfigStruct) string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		c.MySQLUser,
		c.MySQLPass,
		c.MySQLHost,
		c.MySQLPort,
		c.MySQLDB,
	)
}

// assertOpenConnection handles asserting a connection passed into a sql
// function is open, not nil. If nil, we'll create a new connection.
func assertOpenConnection(db *gorm.DB) *gorm.DB {
	var err error

	if db == nil {
		db, err = OpenConnection()
		if err != nil {
			err = err
		}
	}

	return db
}

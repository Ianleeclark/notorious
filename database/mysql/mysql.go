package mysql

import (
	"database/sql"
	"fmt"
	"github.com/GrappigPanda/notorious/config"
	"github.com/GrappigPanda/notorious/database"
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
func InitDB(dbConn *gorm.DB) {
	dbConn = assertOpenConnection(dbConn)

	dbConn.CreateTable(&db.WhiteTorrent{})
	dbConn.CreateTable(&db.Torrent{})
	dbConn.CreateTable(&db.TrackerStats{})
	dbConn.CreateTable(&db.PeerStats{})
}

// GetTorrent retrieves a torrent by its infoHash from the generic torrent
// table in the database. Note: there's also a whitelisted torrent table
// (`WhiteTorrent`).
func GetTorrent(dbConn *gorm.DB, infoHash string) (t *db.Torrent, err error) {
	dbConn = assertOpenConnection(dbConn)

	t = &db.Torrent{}

	dbConn.Where("info_hash = ?", infoHash).Find(&t)

	return
}

// GetWhitelistedTorrent Retrieves a single whitelisted torrent by its infoHash
func GetWhitelistedTorrent(dbConn *gorm.DB, infoHash string) (t *db.WhiteTorrent, err error) {
	dbConn = assertOpenConnection(dbConn)

	t = &db.WhiteTorrent{}

	x := dbConn.Where("info_hash = ?", infoHash).First(&t)
	if x.Error != nil {
		err = x.Error
	}

	return
}

// UpdateStats Handles updating statistics relevant to our tracker.
func UpdateStats(dbConn *gorm.DB, uploaded uint64, downloaded uint64) {
	dbConn = assertOpenConnection(dbConn)

	ts := &db.TrackerStats{}
	dbConn.First(&ts)
	dbConn.Model(&ts).Updates(
		db.TrackerStats{
			Uploaded:   ts.Uploaded + int64(uploaded),
			Downloaded: ts.Downloaded + int64(downloaded),
		})

	return
}

// UpdateTorrentStats Handles updating statistics relevant to our tracker.
func UpdateTorrentStats(dbConn *gorm.DB, seederDelta int64, leecherDelta int64) {
	dbConn = assertOpenConnection(dbConn)

	t := &db.Torrent{}
	dbConn.First(&t)
	dbConn.Model(&t).Updates(
		db.TrackerStats{
			Uploaded:   t.Seeders + seederDelta,
			Downloaded: t.Leechers + leecherDelta,
		})

	return
}

// UpdatePeerStats handles updating peer info like hits per ip, downloaded
// amount, uploaded amounts.
func UpdatePeerStats(dbConn *gorm.DB, uploaded uint64, downloaded uint64, ip string) {
	dbConn = assertOpenConnection(dbConn)

	ps := &db.PeerStats{Ip: ip}
	dbConn.First(&ps)
	dbConn.Model(&ps).UpdateColumn(map[string]interface{}{
		"Uploaded":   ps.Uploaded + int64(uploaded),
		"Downloaded": ps.Downloaded + int64(downloaded),
	})

	return
}

// GetWhitelistedTorrents allows us to retrieve all of the white listed
// torrents. Mostly used for populating the Redis KV storage with all of our
// whitelisted torrents.
func GetWhitelistedTorrents(dbConn *gorm.DB) (x *sql.Rows, err error) {
	dbConn = assertOpenConnection(dbConn)

	x, err = dbConn.Table("white_torrents").Rows()
	if err != nil {
		return
	}

	return
}

// ScrapeTorrent supports the Scrape convention
func ScrapeTorrent(dbConn *gorm.DB, infoHash string) (torrent *db.Torrent) {
	dbConn = assertOpenConnection(dbConn)

	dbConn.Where("info_hash = ?", infoHash).First(&torrent)
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

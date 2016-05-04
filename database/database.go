package db

import (
	"fmt"
	"github.com/GrappigPanda/notorious/config"
	"github.com/jinzhu/gorm"
	// We use a blank import here because I'm afraid of breaking anything
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// formatConnectStrings concatenates the data from the config file into a
// usable MySQL connection string.
func formatConnectString(c config.ConfigStruct) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?parseTime=true",
		c.MySQLUser,
		c.MySQLPass,
		c.MySQLHost,
		c.MySQLPort,
		c.MySQLDB,
	)
}

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
	db.CreateTable(&Torrent{})
	db.CreateTable(&White_Torrent{})
}

// AddWhitelistedTorrent adds a torrent to the whitelist so that they may be
// used by the tracker in the future.
func (t *Torrent) AddWhitelistedTorrent() bool {
	db, err := OpenConnection()
	if err != nil {
		err = err
	}

	db.Create(t)
	return db.NewRecord(t)
}

// GetTorrent retrieves a torrent by its infoHash from the generic torrent
// table in the database. Note: there's also a whitelisted torrent table
// (`white_torrent`).
func GetTorrent(infoHash string) (t *Torrent, err error) {
	db, err := OpenConnection()
	if err != nil {
		err = err
	}
	t = &Torrent{}

	db.Where("info_hash = ?", infoHash).Find(&t)

	return
}

// GetWhitelistedTorrent Retrieves a single whitelisted torrent by its infoHash
func GetWhitelistedTorrent(infoHash string) (t *White_Torrent, err error) {
	db, err := OpenConnection()
	if err != nil {
		err = err
	}
	t = &White_Torrent{}

    x := db.Where("info_hash = ?", infoHash).First(&t)
    if x.Error != nil {
        err = x.Error
    }

	return
}

// GetWhitelistedTorrent allows us to retrieve all of the white listed
// torrents. Mostly used for populating the Redis KV storage with all of our
// whitelisted torrents.
func GetWhitelistedTorrents() (t *White_Torrent, err error) {
	db, err := OpenConnection()
	if err != nil {
		err = err
	}
	t = &White_Torrent{}

    x := db.Where("info_hash = ?", infoHash)
    if x.Error != nil {
        err = x.Error
    }

	return
}


// ScrapeTorrent supports the Scrape convention
func ScrapeTorrent(db *gorm.DB, infoHash string) interface{} {
	var torrent Torrent
	return db.Where("infoHash = ?", infoHash).First(&torrent)
}

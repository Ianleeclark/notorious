package db

import (
	"fmt"
	"github.com/GrappigPanda/notorious/config"
	"github.com/jinzhu/gorm"
    // We use a blank import here because I'm afraid of breaking anything
	_ "github.com/jinzhu/gorm/dialects/mysql"

)

func formatConnectString(c config.ConfigStruct) string {
	return fmt.Sprintf("%s:%s@%s/%s",
		c.MySQLUser,
		c.MySQLPass,
		c.MySQLHost,
		c.MySQLDB,
	)
}

// OpenConnection does as its name dictates and opens a connection to the
// MysqlHost listed in the config
func OpenConnection() (db *gorm.DB, err error) {
	//c := config.LoadConfig()

    db, err = gorm.Open("mysql", "localhost:3306/testdb")
	if err != nil {
        err = fmt.Errorf("Failed to open connection to MySQL: %v", err)
	}

	return
}

func (t *Torrent) AddWhitelistedTorrent() {
    db, err := OpenConnection()
    if err != nil {
        err = err
    }

    db.Create(t)

    return
}

func GetTorrent(infoHash string) (t *Torrent, err error) {
    db, err := OpenConnection()
    if err != nil {
        err = err
    }
    t = &Torrent{}

    db.Where("infoHash = ?", infoHash).First(&t)

    return
}

// ScrapeTorrent supports the Scrape convention
func ScrapeTorrent(db *gorm.DB, infoHash string) interface{} {
	var torrent Torrent
    // TODO(ian): FInish this.
	return db.Where("infoHash = ?", infoHash).Find(&torrent).Value
}

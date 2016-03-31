package db

import (
	"fmt"
	"github.com/GrappigPanda/notorious/config"
	"github.com/jinzhu/gorm"
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

func OpenConnection() *gorm.DB {
	c := config.LoadConfig()

	db, err := gorm.Open("mysql", formatConnectString(c))
	if err != nil {
		panic(err)
	}

	return db
}

func ScrapeTorrent(db *gorm.DB, infoHash string) interface{} {
	var torrent Torrent
	return db.Where("infoHash = ?", infoHash).Find(&torrent).Value
}

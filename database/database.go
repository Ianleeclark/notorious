package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func formatConnectString(c *configStruct) string {
	return fmt.Sprintf("%s:%s@%s/%s",
		c.MySQLUser,
		c.MySQLPass,
		c.MySQLHost,
		c.MySQLDB,
	)
}

func OpenConnection(c *configStruct) *DB {
	db, err := gorm.Open("mysql", formatConnectString(c))
	if err != nil {
		panic(err)
	}

	return db
}

func ScrapeTorrent(db *DB, infoHash string) []string {
	return db.Where("infoHash", infoHash).Find(&Torrent)
}

package db

import (
	"github.com/jinzhu/gorm"
    // We use a blank import here because I'm afraid of breaking anything
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Torrent houses the torrent schema information that we store in the DB.
type Torrent struct {
	gorm.Model
	id         int    `gorm:"AUTO_INCREMENT, unique"`
	infoHash   string `gorm:"varchar(32), not null"`
	name       string `gorm:"not null"`
	Downloaded int    `gorm:"not null"`
	Seeders    int    `gorm:"not null"`
	Leechers   int    `gorm:"not null"`
}

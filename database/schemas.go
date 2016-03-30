package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Torrent struct {
	gorm.Model
	id         int    `gorm:"AUTO_INCREMENT, unique"`
	infoHash   string `gorm:"varchar(32), not null"`
	name       string `gorm:"not null"`
	Downloaded int    `gorm:"not null"`
	Seeders    int    `gorm:"not null"`
	Leechers   int    `gorm:"not null"`
}

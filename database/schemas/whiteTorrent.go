package schemas

import (
	"github.com/jinzhu/gorm"
)

type WhiteTorrent struct {
	id        int    `gorm:"AUTO_INCREMENT, unique, primary_key"`
	InfoHash  string `gorm:"varchar(32), not null"`
	Name      string `gorm:"not null"`
	AddedBy   string `gorm:"varchar(15)"`
	DateAdded int64
}

func (t *WhiteTorrent) AddWhitelistedTorrent(db *gorm.DB) bool {
	db.Create(t)
	return db.NewRecord(t)
}

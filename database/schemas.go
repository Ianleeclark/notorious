package db

// Torrent houses the torrent schema information that we store in the DB.
type Torrent struct {
	id         int    `gorm:"AUTO_INCREMENT, unique, primary_key"`
	InfoHash   string `gorm:"varchar(32), not null"`
	Name       string `gorm:"not null"`
	Downloaded int    `gorm:"not null"`
	Seeders    int64  `gorm:"not null"`
	Leechers   int64  `gorm:"not null"`
	AddedBy    string `gorm:"varchar(15)"`
	DateAdded  int64
}

type WhiteTorrent struct {
	id        int    `gorm:"AUTO_INCREMENT, unique, primary_key"`
	InfoHash  string `gorm:"varchar(32), not null"`
	Name      string `gorm:"not null"`
	AddedBy   string `gorm:"varchar(15)"`
	DateAdded int64
}

type PeerStats struct {
	id         int    `gorm:"AUTO_INCREMENT, unique, primary_key"`
	Downloaded int64  `gorm:"not null"`
	Uploaded   int64  `gorm:"not null"`
	Ip         string `gorm:"varchar(15)"`
}

type TrackerStats struct {
	id         int   `gorm:"AUTO_INCREMENT, unique, primary_key"`
	Downloaded int64 `gorm:"not null"`
	Uploaded   int64 `gorm:"not null"`
}

package schemas

type TrackerStats struct {
	id         int   `gorm:"AUTO_INCREMENT, unique, primary_key"`
	Downloaded int64 `gorm:"not null"`
	Uploaded   int64 `gorm:"not null"`
}

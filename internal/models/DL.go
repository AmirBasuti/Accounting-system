package models

type DL struct {
	ID      uint
	Code    string `gorm:"size:64;not null;unique"`
	Title   string `gorm:"size:64;not null;unique"`
	Version uint   `gorm:"default:1;not null"`
}

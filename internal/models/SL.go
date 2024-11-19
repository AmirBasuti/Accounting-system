package models

type SL struct {
	ID       uint
	Code     string `gorm:"size:64;not null;unique"`
	Title    string `gorm:"size:64;not null;unique"`
	IsDetail bool   `gorm:"not null"`
	Version  uint   `gorm:"default:1;not null"`
}

package models

type VoucherItem struct {
	ID        uint
	VoucherID uint  `gorm:"not null"`
	SLID      uint  `gorm:"not null"`
	DLID      *uint `gorm:""`
	Debit     int   `gorm:"check:debit >= 0"`
	Credit    int   `gorm:"check:credit >= 0"`
}

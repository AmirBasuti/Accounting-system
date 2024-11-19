package models

type Voucher struct {
	ID      uint
	Number  string        `gorm:"size:64;not null;unique"`
	Version uint          `gorm:"default:1;not null"`
	Item    []VoucherItem `gorm:"foreignKey:VoucherID"`
}

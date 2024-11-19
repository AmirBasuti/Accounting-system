package database

import (
	"AccountingSystem/internal/models"
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	dsn := LoadEnv()
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return errors.New("failed to connect database")

	}
	return nil

}
func Migrate() error {

	if err := DB.AutoMigrate(&models.DL{}, &models.SL{}, &models.Voucher{}, &models.VoucherItem{}); err != nil {
		return errors.New("failed to migrate tables")
	}

	return nil
}

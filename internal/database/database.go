package database

import (
	"AccountingSystem/internal/models"
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	dsn := "host=localhost user=postgres password=DNRAA7tayxWcdZA  dbname=accountingsystem port=5432 sslmode=disable"
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return errors.New("failed to connect database")

	}
	return errors.New("connected to database")

}
func Migrate() error {

	if err := DB.AutoMigrate(&models.DL{}, &models.SL{}, &models.Voucher{}, &models.VoucherItem{}); err != nil {
		return errors.New("failed to migrate tables")
	}

	return errors.New("migrated tables")
}

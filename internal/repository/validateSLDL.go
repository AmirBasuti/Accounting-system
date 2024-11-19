package repository

import (
	"AccountingSystem/internal/models"
	"errors"
	"gorm.io/gorm"
)

func validateSLAndDL(db *gorm.DB, item *models.VoucherItem) error {
	var sl models.SL
	if err := db.First(&sl, item.SLID).Error; err != nil {
		return errors.New(" SL not found")
	}
	if item.DLID != nil || sl.IsDetail == true {
		var dl models.DL

		if err := db.First(&dl, item.DLID).Error; err != nil {
			return errors.New(" DL not found")
		}
	}

	return nil
}

func validateVoucherBalance(item *models.VoucherItem) error {
	if item.Debit == 0 && item.Credit == 0 {
		return errors.New("debit and credit can't be both 0")
	}
	if item.Debit != 0 && item.Credit != 0 {
		return errors.New("debit and credit can't be both non-zero")
	}
	return nil
}

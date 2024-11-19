package repository

import (
	"AccountingSystem/internal/models"
	"errors"
	"gorm.io/gorm"
)

// validateSLAndDL checks if the SL and DL fields in a VoucherItem are valid
func validateSLAndDL(db *gorm.DB, item *models.VoucherItem) error {
	// Check if the SL exists
	var sl models.SL
	if err := db.First(&sl, item.SLID).Error; err != nil {
		return errors.New("SL not found")
	}

	// If SL is detail-oriented, check if DL exists
	if item.DLID != nil || sl.IsDetail == true {
		var dl models.DL
		if err := db.First(&dl, item.DLID).Error; err != nil {
			return errors.New("DL not found")
		}
	}

	return nil
}

// validateVoucherBalance checks if either debit or credit is set, but not both
func validateVoucherBalance(item *models.VoucherItem) error {
	// Both debit and credit can't be zero
	if item.Debit == 0 && item.Credit == 0 {
		return errors.New("debit and credit can't both be 0")
	}

	// Both debit and credit can't be non-zero
	if item.Debit != 0 && item.Credit != 0 {
		return errors.New("debit and credit can't both be non-zero")
	}

	return nil
}

// validateVoucherItemCount ensures a voucher has between 2 and 500 items
func validateVoucherItemCount(db *gorm.DB, voucher *models.Voucher) error {
	var count int64

	// Count the items in the voucher
	db.Model(&models.VoucherItem{}).Where("voucher_id = ?", voucher.ID).Count(&count)

	// Check if the count is below the minimum
	if count < 2 {
		return errors.New("voucher must have at least two items")
	}

	// Check if the count exceeds the maximum
	if count > 500 {
		return errors.New("voucher must have at most 500 items")
	}

	return nil
}

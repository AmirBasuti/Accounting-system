package repository

import (
	"AccountingSystem/internal/models"
	"errors"
	"gorm.io/gorm"
)

// VoucherRepo handles database operations for Vouchers
type VoucherRepo struct {
	DB *gorm.DB
}

// Create adds a new Voucher and its items to the database
func (r *VoucherRepo) Create(voucher *models.Voucher, items []*models.VoucherItem) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Save the voucher
		if err := tx.Create(voucher).Error; err != nil {
			return err
		}

		// Save all voucher items
		for _, item := range items {
			item.VoucherID = voucher.ID
			if err := tx.Create(item).Error; err != nil {
				return err
			}
		}

		// Validate the voucher
		if err := validateVoucher(tx, voucher); err != nil {
			return err
		}

		return nil
	})
}

// Delete removes a Voucher and its items from the database
func (r *VoucherRepo) Delete(voucherID, version uint) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var current models.Voucher

		// Check if the voucher exists
		if err := tx.First(&current, voucherID).Error; err != nil {
			return errors.New("voucher not found" + err.Error())
		}

		// Check version for consistency
		if current.Version != version {
			return errors.New("version mismatch")
		}

		// Delete associated items
		if err := tx.Where("voucher_id = ?", voucherID).Delete(&models.VoucherItem{}).Error; err != nil {
			return errors.New("failed to delete voucher items" + err.Error())
		}

		// Delete the voucher
		if err := tx.Delete(&current).Error; err != nil {
			return errors.New("failed to delete voucher" + err.Error())
		}

		return nil
	})
}

// GetByID retrieves a Voucher by its ID, including its items
func (r *VoucherRepo) GetByID(id uint) (*models.Voucher, error) {
	var voucher models.Voucher

	// Use Preload to load related items
	if err := r.DB.Preload("Item").First(&voucher, id).Error; err != nil {
		return nil, errors.New("voucher not found" + err.Error())
	}

	return &voucher, nil
}

// Update modifies a Voucher and handles its related items
func (r *VoucherRepo) Update(voucher *models.Voucher, inserted, updated []*models.VoucherItem, deleted []uint) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var current models.Voucher

		// Check if the voucher exists
		if err := tx.First(&current, voucher.ID).Error; err != nil {
			return errors.New("voucher not found" + err.Error())
		}

		// Check version for consistency
		if current.Version != voucher.Version {
			return errors.New("version mismatch")
		}

		// Insert new items
		for _, item := range inserted {
			item.VoucherID = voucher.ID
			if err := tx.Create(&item).Error; err != nil {
				return errors.New("failed to insert voucher item: " + err.Error())
			}
		}

		// Update existing items
		for _, item := range updated {
			if err := tx.Save(&item).Error; err != nil {
				return errors.New("failed to update voucher item")
			}
		}

		// Delete items
		if len(deleted) > 0 {
			if err := tx.Where("id IN ?", deleted).Delete(&models.VoucherItem{}).Error; err != nil {
				return errors.New("failed to delete voucher items" + err.Error())
			}
		}

		// Validate the voucher
		if err := validateVoucher(tx, voucher); err != nil {
			return errors.New("voucher validation failed: " + err.Error())
		}

		// Update the voucher version and save it
		voucher.Version = current.Version + 1
		if err := tx.Save(voucher).Error; err != nil {
			return errors.New("failed to update voucher" + err.Error())
		}

		return nil
	})
}

// validateVoucher checks the voucher for business rules
func validateVoucher(db *gorm.DB, voucher *models.Voucher) error {
	// Validate item count
	if err := validateVoucherItemCount(db, voucher); err != nil {
		return err
	}

	// Validate SL, DL, and balance for each item
	var items []models.VoucherItem
	db.Model(&models.VoucherItem{}).Where("voucher_id = ?", voucher.ID).Find(&items)

	for _, item := range items {
		if err := validateSLAndDL(db, &item); err != nil {
			return err
		}
		if err := validateVoucherBalance(&item); err != nil {
			return err
		}
	}

	return nil
}

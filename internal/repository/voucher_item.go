package repository

import (
	"AccountingSystem/internal/models"
	"errors"
	"gorm.io/gorm"
)

// VoucherItemRepo handles database operations for VoucherItem
type VoucherItemRepo struct {
	DB *gorm.DB
}

// Delete removes a VoucherItem by its ID
func (r *VoucherItemRepo) Delete(id uint) error {
	return r.DB.Delete(&models.VoucherItem{}, id).Error
}

// GetByID retrieves a VoucherItem by its ID
func (r *VoucherItemRepo) GetByID(id uint) (*models.VoucherItem, error) {
	var item models.VoucherItem

	// Try to find the VoucherItem
	if err := r.DB.First(&item, id).Error; err != nil {
		return nil, errors.New("voucher item not found" + err.Error())
	}

	return &item, nil
}

// Update modifies an existing VoucherItem
// It validates the SL and DL before updating
func (r *VoucherItemRepo) Update(item *models.VoucherItem) error {
	// Validate SL and DL fields
	if err := validateSLAndDL(r.DB, item); err != nil {
		return err
	}

	return r.DB.Save(item).Error
}

// Create adds a new VoucherItem to the database
// It validates the SL and DL before saving
func (r *VoucherItemRepo) Create(item *models.VoucherItem) error {
	// Validate SL and DL fields
	if err := validateSLAndDL(r.DB, item); err != nil {
		return err
	}

	return r.DB.Create(item).Error
}

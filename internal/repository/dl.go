package repository

import (
	"AccountingSystem/internal/models"
	"errors"
	"gorm.io/gorm"
)

// Constant for "detail not found" error message
const (
	ErrDetailNotFound = "detail not found"
)

// DLRepo is the struct for handling DL-related database operations
type DLRepo struct {
	DB *gorm.DB
}

// Create adds a new detail (DL) record to the database
func (r *DLRepo) Create(detail *models.DL) error {
	return r.DB.Create(detail).Error
}

// Update modifies an existing detail (DL) record in the database
// Checks the version to avoid conflicts with other updates
func (r *DLRepo) Update(detail *models.DL) error {
	var current models.DL

	// Find the existing detail by ID
	if err := r.DB.First(&current, detail.ID).Error; err != nil {
		return errors.New(ErrDetailNotFound)
	}

	// Check if the version matches to avoid conflicts
	if current.Version != detail.Version {
		return errors.New("version mismatch: the record has been updated by another process")
	}

	// Increment the version and save the updated record
	detail.Version = current.Version + 1
	return r.DB.Save(detail).Error
}

// Delete removes a detail (DL) record from the database
// It only allows deletion if the detail is not used in any vouchers
func (r *DLRepo) Delete(id, version uint) error {
	var current models.DL

	// Check if the detail exists
	if err := r.DB.First(&current, id).Error; err != nil {
		return errors.New(ErrDetailNotFound)
	}

	// Check if the version matches to avoid conflicts
	if current.Version != version {
		return errors.New("version mismatch: the record has been updated by another process")
	}

	// Check if the detail is used in voucher items
	var count int64
	r.DB.Model(&models.VoucherItem{}).Where("detail_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("detail is used in vouchers and cannot be deleted")
	}

	// Delete the detail
	return r.DB.Delete(current).Error
}

// GetByID finds a detail (DL) record in the database by its ID
func (r *DLRepo) GetByID(id uint) (*models.DL, error) {
	var detail models.DL

	// Try to find the detail
	if err := r.DB.First(&detail, id).Error; err != nil {
		return nil, errors.New(ErrDetailNotFound)
	}

	return &detail, nil
}

package repository

import (
	"AccountingSystem/internal/models"
	"errors"
	"gorm.io/gorm"
)

// SLRepo is a struct to handle SL-related database operations
type SLRepo struct {
	DB *gorm.DB
}

// Create adds a new SL record to the database
func (r *SLRepo) Create(sl *models.SL) error {
	return r.DB.Create(sl).Error
}

// GetByID retrieves an SL record by its ID
func (r *SLRepo) GetByID(id uint) (*models.SL, error) {
	var sl models.SL

	// Try to find the SL by ID
	if err := r.DB.First(&sl, id).Error; err != nil {
		return nil, err
	}

	return &sl, nil
}

// Update modifies an existing SL record in the database
// It ensures the record is not referenced and checks the version
func (r *SLRepo) Update(sl *models.SL) error {
	var current models.SL

	// Find the existing SL record
	if err := r.DB.First(&current, sl.ID).Error; err != nil {
		return err
	}

	// Check version to avoid conflicts
	if current.Version != sl.Version {
		return errors.New("version mismatch: the record has been updated by another process")
	}

	// Check if the SL is used in voucher items
	var count int64
	r.DB.Model(&models.VoucherItem{}).Where("sl_id = ?", sl.ID).Count(&count)
	if count > 0 {
		return errors.New("SL is referenced and cannot be updated")
	}

	// Update the version and save the record
	sl.Version = current.Version + 1

	return r.DB.Save(sl).Error
}

// Delete removes an SL record from the database
// It ensures the SL is not used in voucher items before deleting
func (r *SLRepo) Delete(id, version uint) error {
	var sl models.SL

	// Check if the SL exists
	if err := r.DB.First(&sl, id).Error; err != nil {
		return errors.New("SL not found")
	}

	// Check version to avoid conflicts
	if sl.Version != version {
		return errors.New("version mismatch: the record has been updated by another process")
	}

	// Check if the SL is used in voucher items
	var count int64
	r.DB.Model(&models.VoucherItem{}).Where("sl_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("SL is referenced and cannot be deleted")
	}

	// Delete the SL record
	return r.DB.Delete(&sl).Error
}

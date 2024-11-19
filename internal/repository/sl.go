package repository

import (
	"AccountingSystem/internal/models"
	"errors"
	"gorm.io/gorm"
)

type SLRepo struct {
	DB *gorm.DB
}

func (r *SLRepo) Create(sl *models.SL) error {
	return r.DB.Create(sl).Error
}

func (r *SLRepo) GetByID(id uint) (*models.SL, error) {
	var sl models.SL

	if err := r.DB.First(&sl, id).Error; err != nil {
		return nil, err
	}

	return &sl, nil
}

func (r *SLRepo) Update(sl *models.SL) error {
	var current models.SL

	if err := r.DB.First(&current, sl.ID).Error; err != nil {
		return err
	}
	if current.Version != sl.Version {
		return errors.New("version mismatch: the record has been updated by another process")
	}

	var count int64
	r.DB.Model(&models.VoucherItem{}).Where("sl_id = ?", sl.ID).Count(&count)
	if count > 0 {
		return errors.New("SL is referenced and cannot be updated")
	}

	sl.Version = current.Version + 1

	return r.DB.Save(sl).Error
}

func (r *SLRepo) Delete(id, version uint) error {
	var sl models.SL

	if err := r.DB.First(&sl, id).Error; err != nil {
		return errors.New("detail not found")
	}
	if sl.Version != version {
		return errors.New("version mismatch: the record has been updated by another process")
	}

	var count int64
	r.DB.Model(&models.VoucherItem{}).Where("sl_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("SL is referenced and cannot be deleted")
	}

	return r.DB.Delete(&sl).Error
}

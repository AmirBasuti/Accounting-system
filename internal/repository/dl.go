package repository

import (
	"AccountingSystem/internal/models"
	"errors"
	"gorm.io/gorm"
)

type DLRepo struct {
	DB *gorm.DB
}

func (r *DLRepo) Create(detail *models.DL) error {

	return r.DB.Create(detail).Error
}

func (r *DLRepo) Update(detail *models.DL) error {
	var current models.DL

	if err := r.DB.First(&current, detail.ID).Error; err != nil {
		return errors.New("detail not found")
	}

	if current.Version != detail.Version {
		return errors.New("version mismatch: the record has been updated by another process")
	}

	detail.Version = current.Version + 1
	return r.DB.Save(detail).Error

}

func (r *DLRepo) Delete(id, version uint) error {
	var current models.DL

	if err := r.DB.First(&current, id).Error; err != nil {
		return errors.New("detail not found")
	}
	if current.Version != version {
		return errors.New("version mismatch: the record has been updated by another process")
	}

	var count int64
	r.DB.Model(&models.VoucherItem{}).Where("detail_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("detail is used in vouchers and cannot be deleted")
	}

	return r.DB.Delete(current).Error
}

func (r *DLRepo) GetByID(id uint) (*models.DL, error) {
	var detail models.DL

	if err := r.DB.First(&detail, id).Error; err != nil {
		return nil, errors.New("detail not found")
	}
	return &detail, nil

}

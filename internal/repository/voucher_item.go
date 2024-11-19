package repository

import (
	"AccountingSystem/internal/models"
	"errors"
	"gorm.io/gorm"
)

type VoucherItemRepo struct {
	DB *gorm.DB
}

func (r *VoucherItemRepo) Delete(id uint) error {
	return r.DB.Delete(&models.VoucherItem{}, id).Error
}

func (r *VoucherItemRepo) GetByID(id uint) (*models.VoucherItem, error) {
	var item models.VoucherItem
	if err := r.DB.First(&item, id).Error; err != nil {
		return nil, errors.New("voucher item not found")
	}
	return &item, nil
}

func (r *VoucherItemRepo) Update(item *models.VoucherItem) error {
	if err := validateSLAndDL(r.DB, item); err != nil {
		return err
	}
	return r.DB.Save(item).Error
}
func (r *VoucherItemRepo) Create(item *models.VoucherItem) error {
	if err := validateSLAndDL(r.DB, item); err != nil {
		return err
	}

	return r.DB.Create(item).Error
}

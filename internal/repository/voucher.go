package repository

import (
	"AccountingSystem/internal/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type VoucherRepo struct {
	DB *gorm.DB
}

func (r *VoucherRepo) Create(voucher *models.Voucher, items []*models.VoucherItem) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(voucher).Error; err != nil {
			return err
		}
		for _, item := range items {
			item.VoucherID = voucher.ID
			//voucher.Item = append(voucher.Item, *item)
			if err := tx.Create(item).Error; err != nil {
				return err
			}
		}
		if err := validateVoucher(tx, voucher); err != nil {
			return err
		}

		return nil
	})
}
func (r *VoucherRepo) Delete(voucherID, version uint) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var current models.Voucher
		if err := tx.First(&current, voucherID).Error; err != nil {
			return errors.New("voucher not found")
		}
		if current.Version != version {
			return errors.New("version mismatch")
		}
		if err := tx.Where("voucher_id = ?", voucherID).Delete(&models.VoucherItem{}).Error; err != nil {
			return errors.New("failed to delete voucher items")
		}
		if err := tx.Delete(&current).Error; err != nil {
			return errors.New("failed to delete voucher")
		}
		return nil
	})
}
func (r *VoucherRepo) GetByID(id uint) (*models.Voucher, error) {
	var voucher models.Voucher
	if err := r.DB.Preload("Item").First(&voucher, id).Error; err != nil {
		return nil, errors.New("voucher not found")
	}
	return &voucher, nil
}
func (r *VoucherRepo) Update(voucher *models.Voucher, inserted, updated []*models.VoucherItem, deleted []uint) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var current models.Voucher
		if err := tx.First(&current, voucher.ID).Error; err != nil {
			return errors.New("voucher not found")
		}
		if current.Version != voucher.Version {
			fmt.Println(current.Version, voucher.Version)
			return errors.New("version mismatch")
		}
		for _, item := range inserted {
			item.VoucherID = voucher.ID
			if err := tx.Create(&item).Error; err != nil {
				return errors.New("failed to insert voucher item" + err.Error())
			}
		}
		for _, item := range updated {
			if err := tx.Save(&item).Error; err != nil {
				return errors.New("failed to update voucher item")
			}
		}
		if len(deleted) > 0 {
			if err := tx.Where("id IN ?", deleted).Delete(&models.VoucherItem{}).Error; err != nil {
				return errors.New("failed to delete voucher items")
			}
		}
		if err := validateVoucher(tx, voucher); err != nil {
			return errors.New("voucher validation failed" + err.Error())
		}
		voucher.Version = current.Version + 1
		if err := tx.Save(voucher).Error; err != nil {
			return errors.New("failed to update voucher")
		}

		return nil
	})
}

func validateVoucherItemCount(db *gorm.DB, voucher *models.Voucher) error {
	var count int64
	db.Model(&models.VoucherItem{}).Where("voucher_id = ?", voucher.ID).Count(&count)

	if count < 2 {
		return errors.New("voucher must have at least two item")
	}
	if count > 500 {
		return errors.New("voucher must have at most 500 items")
	}
	return nil
}

//	func validateVoucherBalance(items []models.VoucherItem) error {
//		var debit, credit int32
//		for _, item := range items {
//			debit += item.Debit
//			credit += item.Credit
//		}
//		if debit != credit {
//			return errors.New("voucher must be balanced")
//		}
//		return nil
//	}
//
//	func validateSLAndDL(db *gorm.DB, items []models.VoucherItem) error {
//		for _, v := range items {
//			var sl models.SL
//			if err := db.First(&sl, v.SLID).Error; err != nil {
//				return errors.New("SL not found")
//			}
//			if v.DLID != nil {
//				var dl models.DL
//				if err := db.First(&dl, *v.DLID).Error; err != nil {
//					return errors.New("DL not found")
//				}
//			}
//		}
//		return nil
//	}
func validateVoucher(db *gorm.DB, voucher *models.Voucher) error {
	if err := validateVoucherItemCount(db, voucher); err != nil {
		return err
	}
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

package repository_test

import (
	"AccountingSystem/internal/database"
	"AccountingSystem/internal/models"
	"AccountingSystem/internal/repository"
	"testing"
)

func TestSLAndDLReferences(t *testing.T) {
	slRepo := &repository.SLRepo{DB: database.DB}
	dlRepo := &repository.DLRepo{DB: database.DB}
	voucherRepo := &repository.VoucherRepo{DB: database.DB}
	voucherItemRepo := &repository.VoucherItemRepo{DB: database.DB}

	// Create SL and DL
	code, title := repository.GenerateRandomCodeAndTitle()
	sl := &models.SL{Code: code, Title: title, IsDetail: false}
	if err := slRepo.Create(sl); err != nil {
		t.Fatalf("failed to create SL: %v", err)
	}
	code, title = repository.GenerateRandomCodeAndTitle()
	dl := &models.DL{Code: code, Title: title}
	if err := dlRepo.Create(dl); err != nil {
		t.Fatalf("failed to create DL: %v", err)
	}
	code, title = repository.GenerateRandomCodeAndTitle()
	sl2 := &models.SL{Code: code, Title: title, IsDetail: false}
	if err := slRepo.Create(sl2); err != nil {
		t.Fatalf("failed to create SL: %v", err)
	}
	code, title = repository.GenerateRandomCodeAndTitle()
	sl3 := &models.SL{Code: code, Title: title, IsDetail: true}
	if err := slRepo.Create(sl3); err != nil {
		t.Fatalf("failed to create SL: %v", err)
	}
	// Create Voucher and VoucherItem
	voucher_number := repository.GenerateRandomVoucherNumber()
	voucher := &models.Voucher{Number: voucher_number}

	voucherItem := &models.VoucherItem{SLID: sl.ID, Debit: 100, Credit: 0}
	voucherItem2 := &models.VoucherItem{SLID: sl2.ID, Debit: 100, Credit: 0}
	voucherItem3 := &models.VoucherItem{SLID: sl3.ID, Debit: 100, Credit: 0}
	if err := voucherRepo.Create(voucher, []*models.VoucherItem{voucherItem, voucherItem2}); err != nil {
		t.Fatalf("failed to create Voucher: %v", err)
	}

	// Test SL and DL should not be deleted if referenced in VoucherItem
	if err := slRepo.Delete(sl.ID, sl.Version); err == nil {
		t.Fatalf("expected error when deleting referenced SL, got nil")
	}
	//if err := dlRepo.Delete(dl.ID, dl.Version); err == nil {
	//	t.Fatalf("expected error when deleting referenced DL, got nil")
	//}

	// Test SL should not be updated if referenced in VoucherItem
	sl.Title = "Updated Title"
	if err := slRepo.Update(sl); err == nil {
		t.Fatalf("expected error when updating referenced SL, got nil")
	}

	// Test DL must exist in the database if IsDetail is true

	if err := voucherRepo.Update(voucher, []*models.VoucherItem{voucherItem3}, []*models.VoucherItem{}, nil); err == nil {
		t.Fatalf("expected error when updating referenced SL becuse Dl ref is empity, but got nil")
	}
	voucherItem3.DLID = &dl.ID
	if err := voucherRepo.Update(voucher, []*models.VoucherItem{voucherItem3}, []*models.VoucherItem{}, nil); err != nil {
		t.Fatalf("failed to update VoucherItem: %v", err)

	}

	//voucher.DLID = &dl.ID
	//if err := slRepo.Update(sl); err != nil {
	//	t.Fatalf("failed to update SL with DL reference: %v", err)
	//}

	// Cleanup
	if err := voucherItemRepo.Delete(voucherItem.ID); err != nil {
		t.Errorf("failed to delete VoucherItem: %v", err)
	}
	if err := voucherItemRepo.Delete(voucherItem2.ID); err != nil {
		t.Errorf("failed to delete VoucherItem: %v", err)
	}
	if err := voucherItemRepo.Delete(voucherItem3.ID); err != nil {
		t.Errorf("failed to delete VoucherItem: %v", err)
	}

	if err := voucherRepo.Delete(voucher.ID, voucher.Version); err != nil {
		t.Errorf("failed to delete Voucher: %v", err)
	}
	if err := dlRepo.Delete(dl.ID, dl.Version); err != nil {
		t.Errorf("failed to delete DL: %v", err)
	}
	if err := slRepo.Delete(sl.ID, sl.Version); err != nil {
		t.Errorf("failed to delete SL: %v", err)
	}
	if err := slRepo.Delete(sl2.ID, sl2.Version); err != nil {
		t.Errorf("failed to delete SL: %v", err)
	}
	if err := slRepo.Delete(sl3.ID, sl3.Version); err != nil {
		t.Errorf("failed to delete SL: %v", err)
	}

}

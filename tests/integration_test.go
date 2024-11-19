// Package repository_test contains integration tests for the repository package.
// These tests ensure that the SL, DL, Voucher, and VoucherItem repositories
// correctly handle creation, updating, and deletion of records, especially
// when there are references between them.

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
	voucherNumber := repository.GenerateRandomVoucherNumber()
	voucher := &models.Voucher{Number: voucherNumber}

	voucherItem := &models.VoucherItem{SLID: sl.ID, DLID: nil, Debit: 100, Credit: 0}
	voucherItem2 := &models.VoucherItem{SLID: sl2.ID, DLID: nil, Debit: 100, Credit: 0}
	voucherItem3 := &models.VoucherItem{SLID: sl3.ID, DLID: nil, Debit: 100, Credit: 0}
	if err := voucherRepo.Create(voucher, []*models.VoucherItem{voucherItem, voucherItem2}); err != nil {
		// Attempting to delete SL should fail because it is referenced in a VoucherItem
		if err := slRepo.Delete(sl.ID, sl.Version); err == nil {
			t.Fatalf("expected error when deleting referenced SL, got nil")
		}
		// Test SL and DL should not be deleted if referenced in VoucherItem
		if err := slRepo.Delete(sl.ID, sl.Version); err == nil {
			t.Fatalf("expected error when deleting referenced SL, got nil")
		}

		// Test SL should not be updated if referenced in VoucherItem
		// This update is expected to fail because the SL is referenced in a VoucherItem.
		sl.Title = "Updated Title"
		if err := slRepo.Update(sl); err == nil {
			t.Fatalf("expected error when updating referenced SL, got nil")
		}
	}

	// Test DL must exist in the database if IsDetail is true

	if err := voucherRepo.Update(voucher, []*models.VoucherItem{voucherItem3}, []*models.VoucherItem{}, nil); err == nil {
		t.Fatalf("expected error when updating referenced SL because DL reference is empty, but got nil")
	}
	// Set DLID to dl.ID because SL with IsDetail true must have a corresponding DL reference
	voucherItem3.DLID = &dl.ID
	if err := voucherRepo.Update(voucher, []*models.VoucherItem{voucherItem3}, []*models.VoucherItem{}, nil); err != nil {
		t.Fatalf("failed to update VoucherItem: %v", err)

	}

	// Cleanup the created records
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
		// The DL is no longer referenced by any VoucherItem, so it can be safely deleted.
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

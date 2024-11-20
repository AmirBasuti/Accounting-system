package repository_test

import (
	"AccountingSystem/internal/database"
	"AccountingSystem/internal/models"
	"AccountingSystem/internal/repository"
	"testing"
)

func TestVoucherItemRepo_Create(t *testing.T) {
	slRepo := &repository.SLRepo{DB: database.DB}
	voucherRepo := &repository.VoucherRepo{DB: database.DB}
	voucherItemRepo := &repository.VoucherItemRepo{DB: database.DB}
	createdVoucherItems := make(map[uint]struct{})
	createdVouchers := make(map[uint]uint)

	// Dynamically create an SL (Sub Ledger)
	slCode, slTitle := repository.GenerateRandomCodeAndTitle()
	sl := &models.SL{
		Code:  slCode,
		Title: slTitle,
	}
	if err := slRepo.Create(sl); err != nil {
		t.Fatalf("failed to create SL: %v", err)
	}

	sl2Code, sl2Title := repository.GenerateRandomCodeAndTitle()
	sl2 := &models.SL{
		Code:  sl2Code,
		Title: sl2Title,
	}
	if err := slRepo.Create(sl2); err != nil {
		t.Fatalf("failed to create SL: %v", err)
	}

	// Use dynamically created SLID and VoucherID
	voucherItem := &models.VoucherItem{
		SLID:   sl.ID, // From the dynamically created SL
		DLID:   nil,   // Assuming DLID is nullable
		Debit:  100,
		Credit: 0,
	}

	voucherItem2 := &models.VoucherItem{
		SLID:   sl2.ID, // From the dynamically created SL
		DLID:   nil,    // Assuming DLID is nullable
		Debit:  100,
		Credit: 0,
	}

	// Dynamically create a Voucher
	voucherNumber := repository.GenerateRandomVoucherNumber()
	voucher := &models.Voucher{
		Number: voucherNumber,
	}
	if err := voucherRepo.Create(voucher, []*models.VoucherItem{voucherItem2, voucherItem}); err != nil {
		t.Fatalf("failed to create Voucher: %v", err)
	}
	createdVouchers[voucher.ID] = voucher.Version
	createdVoucherItems[voucherItem.ID] = struct{}{}
	createdVoucherItems[voucherItem2.ID] = struct{}{}

	if voucherItem.ID == 0 {
		t.Fatalf("expected: VoucherItem ID > 0, got: %d", voucherItem.ID)
	}

	// Additional assertions for valid creation
	if voucherItem.VoucherID != voucher.ID {
		t.Errorf("expected: VoucherID = %d, got: %d", voucher.ID, voucherItem.VoucherID)
	}
	if voucherItem.SLID != sl.ID {
		t.Errorf("expected: SLID = %d, got: %d", sl.ID, voucherItem.SLID)
	}

	// Cleanup: Delete created records
	for id := range createdVoucherItems {
		if err := voucherItemRepo.Delete(id); err != nil {
			t.Errorf("failed to delete VoucherItem: %v", err)
		}
	}
	for id, version := range createdVouchers {
		if err := voucherRepo.Delete(id, version); err != nil {
			t.Errorf("failed to delete Voucher: %v", err)
		}
	}
	if err := slRepo.Delete(sl.ID, sl.Version); err != nil {
		t.Errorf("failed to delete SL: %v", err)
	}
	if err := slRepo.Delete(sl2.ID, sl2.Version); err != nil {
		t.Errorf("failed to delete SL: %v", err)
	}
}

func TestVoucherItemRepo_Create_WithDL(t *testing.T) {
	slRepo := &repository.SLRepo{DB: database.DB}
	dlRepo := &repository.DLRepo{DB: database.DB}
	voucherRepo := &repository.VoucherRepo{DB: database.DB}
	voucherItemRepo := &repository.VoucherItemRepo{DB: database.DB}
	createdVoucherItems := make(map[uint]struct{})
	createdVouchers := make(map[uint]uint)

	createdDLs := make(map[uint]uint)
	// Maps to track created records

	// Dynamically create an SL (Sub Ledger)
	slCode, slTitle := repository.GenerateRandomCodeAndTitle()
	sl := &models.SL{
		Code:  slCode,
		Title: slTitle,
	}
	if err := slRepo.Create(sl); err != nil {
		t.Fatalf("failed to create SL: %v", err)
	}

	// Dynamically create a DL (Detail Ledger)
	dlCode, dlTitle := repository.GenerateRandomCodeAndTitle()
	dl := &models.DL{
		Code:  dlCode,
		Title: dlTitle,
	}
	if err := dlRepo.Create(dl); err != nil {
		t.Fatalf("failed to create DL: %v", err)
	}
	createdDLs[dl.ID] = dl.Version

	// Use dynamically created SLID and DLID
	voucherItem1 := &models.VoucherItem{
		SLID:   sl.ID,  // From the dynamically created SL
		DLID:   &dl.ID, // From the dynamically created DL
		Debit:  100,
		Credit: 0,
	}

	voucherItem2 := &models.VoucherItem{
		SLID:   sl.ID,  // From the dynamically created SL
		DLID:   &dl.ID, // From the dynamically created DL
		Debit:  200,
		Credit: 0,
	}

	// Dynamically create a Voucher
	voucherNumber := repository.GenerateRandomVoucherNumber()
	voucher := &models.Voucher{
		Number: voucherNumber,
	}
	if err := voucherRepo.Create(voucher, []*models.VoucherItem{voucherItem1, voucherItem2}); err != nil {
		t.Fatalf("failed to create Voucher: %v", err)
	}
	createdVouchers[voucher.ID] = voucher.Version
	createdVoucherItems[voucherItem1.ID] = struct{}{}
	createdVoucherItems[voucherItem2.ID] = struct{}{}

	if voucherItem1.ID == 0 {
		t.Fatalf("expected: VoucherItem ID > 0, got: %d", voucherItem1.ID)
	}

	// Additional assertions for valid creation
	if voucherItem1.VoucherID != voucher.ID {
		t.Errorf("expected: VoucherID = %d, got: %d", voucher.ID, voucherItem1.VoucherID)
	}
	if voucherItem1.SLID != sl.ID {
		t.Errorf("expected: SLID = %d, got: %d", sl.ID, voucherItem1.SLID)
	}
	if voucherItem1.DLID == nil || *voucherItem1.DLID != dl.ID {
		t.Errorf("expected: DLID = %d, got: %v", dl.ID, voucherItem1.DLID)
	}

	// Cleanup: Delete created records
	for id := range createdVoucherItems {
		if err := voucherItemRepo.Delete(id); err != nil {
			t.Errorf("failed to delete VoucherItem: %v", err)
		}
	}
	for id, version := range createdVouchers {
		if err := voucherRepo.Delete(id, version); err != nil {
			t.Errorf("failed to delete Voucher: %v", err)
		}
	}
	for id, version := range createdDLs {
		if err := dlRepo.Delete(id, version); err != nil {
			t.Errorf("failed to delete DL: %v", err)
		}
	}
	if err := slRepo.Delete(sl.ID, sl.Version); err != nil {
		t.Errorf("failed to delete SL: %v", err)

	}
}

func TestVoucherRepo_Update(t *testing.T) {
	slRepo := &repository.SLRepo{DB: database.DB}
	dlRepo := &repository.DLRepo{DB: database.DB}
	voucherRepo := &repository.VoucherRepo{DB: database.DB}
	voucherItemRepo := &repository.VoucherItemRepo{DB: database.DB}
	slCode, slTitle := repository.GenerateRandomCodeAndTitle()
	sl := &models.SL{
		Code:  slCode,
		Title: slTitle,
	}
	if err := slRepo.Create(sl); err != nil {
		t.Fatalf("failed to create SL: %v", err)
	}

	// Step 2: Create a DL (Detail Ledger)
	dlCode, dlTitle := repository.GenerateRandomCodeAndTitle()
	dl := &models.DL{
		Code:  dlCode,
		Title: dlTitle,
	}
	if err := dlRepo.Create(dl); err != nil {
		t.Fatalf("failed to create DL: %v", err)
	}

	// Step 3: Create a Voucher with two VoucherItems
	voucherNumber := repository.GenerateRandomVoucherNumber()
	voucher := &models.Voucher{
		Number: voucherNumber,
	}
	voucherItem1 := &models.VoucherItem{
		SLID:   sl.ID,
		DLID:   &dl.ID,
		Debit:  100,
		Credit: 0,
	}
	voucherItem2 := &models.VoucherItem{
		SLID:   sl.ID,
		DLID:   &dl.ID,
		Debit:  200,
		Credit: 0,
	}
	if err := voucherRepo.Create(voucher, []*models.VoucherItem{voucherItem1, voucherItem2}); err != nil {
		t.Fatalf("failed to create Voucher: %v", err)
	}

	// Step 4: Update the Voucher and its items
	voucher.Number = "UpdatedNumber"
	voucherItem1.Debit = 150
	// Temporarily set credit to 50 to simulate an invalid update scenario
	voucherItem2.Credit = 50
	if err := voucherRepo.Update(voucher, []*models.VoucherItem{}, []*models.VoucherItem{voucherItem2, voucherItem1}, nil); err == nil {
		t.Fatalf("expected error due to wrong credit, got: nil")
	}
	// Reset credit to 0 for a valid update
	voucherItem2.Credit = 0
	if err := voucherRepo.Update(voucher, []*models.VoucherItem{}, []*models.VoucherItem{voucherItem2, voucherItem1}, nil); err != nil {
		t.Fatalf("failed to update Voucher: %v", err)
	}
	// Step 5: Validate the updated data
	updatedVoucher, err := voucherRepo.GetByID(voucher.ID)
	if err != nil {
		t.Fatalf("failed to retrieve updated Voucher: %v", err)
	}
	if updatedVoucher.Number != "UpdatedNumber" {
		t.Errorf("expected: Voucher Number = UpdatedNumber, got: %s", updatedVoucher.Number)
	}
	if updatedVoucher.Version != 2 {
		t.Errorf("expected: Voucher Version = 2, got: %d", updatedVoucher.Version)
	}

	updatedItem1, err := voucherItemRepo.GetByID(voucherItem1.ID)
	if err != nil {
		t.Fatalf("failed to retrieve updated VoucherItem1: %v", err)
	}
	if updatedItem1.Debit != 150 {
		t.Errorf("expected: VoucherItem1 Debit = 150, got: %d", updatedItem1.Debit)
	}

	updatedItem2, err := voucherItemRepo.GetByID(voucherItem2.ID)
	if err != nil {
		t.Fatalf("failed to retrieve updated VoucherItem2: %v", err)
	}
	if updatedItem2.Credit != 0 {
		t.Errorf("expected: VoucherItem2 Credit = 0, got: %d", updatedItem2.Credit)
	}

	// Step 6: Cleanup created records
	if err := voucherItemRepo.Delete(voucherItem1.ID); err != nil {
		t.Errorf("failed to delete VoucherItem1: %v", err)
	}
	if err := voucherItemRepo.Delete(voucherItem2.ID); err != nil {
		t.Errorf("failed to delete VoucherItem2: %v", err)
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
}

func TestVoucherRepo_GetByID(t *testing.T) {
	slRepo := &repository.SLRepo{DB: database.DB}
	dlRepo := &repository.DLRepo{DB: database.DB}
	voucherRepo := &repository.VoucherRepo{DB: database.DB}
	voucherItemRepo := &repository.VoucherItemRepo{DB: database.DB}

	// Step 1: Create an SL (Sub Ledger)
	slCode, slTitle := repository.GenerateRandomCodeAndTitle()
	sl := &models.SL{
		Code:  slCode,
		Title: slTitle,
	}
	if err := slRepo.Create(sl); err != nil {
		t.Fatalf("failed to create SL: %v", err)
	}

	// Step 2: Create a DL (Detail Ledger)
	dlCode, dlTitle := repository.GenerateRandomCodeAndTitle()
	dl := &models.DL{
		Code:  dlCode,
		Title: dlTitle,
	}
	if err := dlRepo.Create(dl); err != nil {
		t.Fatalf("failed to create DL: %v", err)
	}

	// Step 3: Create a Voucher with VoucherItems
	voucherNumber := repository.GenerateRandomVoucherNumber()
	voucher := &models.Voucher{
		Number: voucherNumber,
	}
	voucherItem1 := &models.VoucherItem{
		SLID:   sl.ID,
		DLID:   &dl.ID,
		Debit:  100,
		Credit: 0,
	}
	voucherItem2 := &models.VoucherItem{
		SLID:   sl.ID,
		DLID:   &dl.ID,
		Debit:  200,
		Credit: 0,
	}
	if err := voucherRepo.Create(voucher, []*models.VoucherItem{voucherItem1, voucherItem2}); err != nil {
		t.Fatalf("failed to create Voucher: %v", err)
	}

	// Step 4: Retrieve the Voucher by ID
	retrievedVoucher, err := voucherRepo.GetByID(voucher.ID)
	if err != nil {
		t.Fatalf("failed to retrieve Voucher: %v", err)
	}

	// Step 5: Validate the retrieved data
	if retrievedVoucher.ID != voucher.ID {
		t.Errorf("expected: Voucher ID = %d, got: %d", voucher.ID, retrievedVoucher.ID)
	}
	if retrievedVoucher.Number != voucher.Number {
		t.Errorf("expected: Voucher Number = %s, got: %s", voucher.Number, retrievedVoucher.Number)
	}
	if len(retrievedVoucher.Item) != 2 {
		t.Errorf("expected: 2 VoucherItems, got: %d", len(retrievedVoucher.Item))
	}

	// Step 6: Cleanup created records
	if err := voucherItemRepo.Delete(voucherItem1.ID); err != nil {
		t.Errorf("failed to delete VoucherItem1: %v", err)
	}
	if err := voucherItemRepo.Delete(voucherItem2.ID); err != nil {
		t.Errorf("failed to delete VoucherItem2: %v", err)
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
}

func TestVoucherRepo_Delete(t *testing.T) {
	slRepo := &repository.SLRepo{DB: database.DB}
	dlRepo := &repository.DLRepo{DB: database.DB}
	voucherRepo := &repository.VoucherRepo{DB: database.DB}
	voucherItemRepo := &repository.VoucherItemRepo{DB: database.DB}

	// Step 1: Create an SL (Sub Ledger)
	slCode, slTitle := repository.GenerateRandomCodeAndTitle()
	sl := &models.SL{
		Code:  slCode,
		Title: slTitle,
	}
	if err := slRepo.Create(sl); err != nil {
		t.Fatalf("failed to create SL: %v", err)
	}

	// Step 2: Create a DL (Detail Ledger)
	dlCode, dlTitle := repository.GenerateRandomCodeAndTitle()
	dl := &models.DL{
		Code:  dlCode,
		Title: dlTitle,
	}
	if err := dlRepo.Create(dl); err != nil {
		t.Fatalf("failed to create DL: %v", err)
	}

	// Step 3: Create a Voucher with VoucherItems
	voucherNumber := repository.GenerateRandomVoucherNumber()
	voucher := &models.Voucher{
		Number: voucherNumber,
	}
	voucherItem1 := &models.VoucherItem{
		SLID:   sl.ID,
		DLID:   &dl.ID,
		Debit:  100,
		Credit: 0,
	}
	voucherItem2 := &models.VoucherItem{
		SLID:   sl.ID,
		DLID:   &dl.ID,
		Debit:  200,
		Credit: 0,
	}
	if err := voucherRepo.Create(voucher, []*models.VoucherItem{voucherItem1, voucherItem2}); err != nil {
		t.Fatalf("failed to create Voucher: %v", err)
	}

	// Step 4: Delete the Voucher by ID and Version
	if err := voucherRepo.Delete(voucher.ID, voucher.Version); err != nil {
		t.Fatalf("failed to delete Voucher: %v", err)
	}

	// Step 5: Validate the deletion
	if _, err := voucherRepo.GetByID(voucher.ID); err == nil {
		t.Errorf("expected: voucher not found, got: nil")
	}

	// Step 6: Cleanup created records
	if err := voucherItemRepo.Delete(voucherItem1.ID); err != nil {
		t.Errorf("failed to delete VoucherItem1: %v", err)
	}
	if _, err := voucherItemRepo.GetByID(voucherItem1.ID); err == nil {
		t.Errorf("expected: voucher item not found, got: nil")

	}
	if err := voucherItemRepo.Delete(voucherItem2.ID); err != nil {
		t.Errorf("failed to delete VoucherItem2: %v", err)
	}
	if _, err := voucherItemRepo.GetByID(voucherItem2.ID); err == nil {
		t.Errorf("expected: voucher item not found, got: nil")
	}
	if err := dlRepo.Delete(dl.ID, dl.Version); err != nil {
		t.Errorf("failed to delete DL: %v", err)
	}
	if err := slRepo.Delete(sl.ID, sl.Version); err != nil {
		t.Errorf("failed to delete SL: %v", err)
	}
}

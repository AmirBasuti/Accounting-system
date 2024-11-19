package main

import (
	"AccountingSystem/internal/database"
)

func main() {
	database.Connect()
	database.Migrate()
	//repo := &repository.VoucherItemRepo{DB: database.DB}
	//slRepo := &repository.SLRepo{DB: database.DB}
	//dlRapo := &repository.DLRepo{DB: database.DB}
	//voucherRepo := &repository.VoucherRepo{DB: database.DB}
	//
	////// Dynamically create an SL (Sub Ledger)
	//slCode, slTitle := repository.GenerateRandomCodeAndTitle()
	//sl := &models.SL{
	//	Code:     slCode,
	//	Title:    slTitle,
	//	IsDetail: true,
	//}
	//if err := slRepo.Create(sl); err != nil {
	//	fmt.Println("failed to create SL: %v", err)
	//}
	////
	//sl2Code, sl2Title := repository.GenerateRandomCodeAndTitle()
	//sl2 := &models.SL{
	//	Code:  sl2Code,
	//	Title: sl2Title,
	//}
	//if err := slRepo.Create(sl2); err != nil {
	//	fmt.Println("failed to create SL: %v", err)
	//}
	//dlcode, dltitle := repository.GenerateRandomCodeAndTitle()
	//dl := &models.DL{
	//	Code:  dlcode,
	//	Title: dltitle,
	//}
	//if err := dlRapo.Create(dl); err != nil {
	//	fmt.Println("failed to create SL: %v", err)
	//
	//}
	////// Use dynamically created SLID and VoucherID
	//voucherItem := &models.VoucherItem{
	//	SLID:   sl.ID, // From the dynamically created SL
	//	DLID:   nil,   // Assuming DLID is nullable
	//	Debit:  100,
	//	Credit: 0,
	//}
	//voucherItem.DLID = &dl.ID
	////// Valid create test
	//////if err := repo.Create(voucherItem); err != nil {
	//////	fmt.Println("failed to create VoucherItem: %v", err)
	//////}
	//voucherItem2 := &models.VoucherItem{
	//	SLID:   sl2.ID, // From the dynamically created SL
	//	DLID:   nil,    // Assuming DLID is nullable
	//	Debit:  100,
	//	Credit: 0,
	//}
	////
	////// Valid create test
	//////if err := repo.Create(voucherItem2); err != nil {
	//////	fmt.Println("failed to create VoucherItem: %v", err)
	//////}
	////
	////// Dynamically create a Voucher
	//voucherNumber := repository.GenerateRandomVoucherNumber()
	//voucher := &models.Voucher{
	//	Number: voucherNumber,
	//}
	//if err := voucherRepo.Create(voucher, []*models.VoucherItem{voucherItem2, voucherItem}); err != nil {
	//	fmt.Println("failed to create Voucher: %v", err)
	//}

}

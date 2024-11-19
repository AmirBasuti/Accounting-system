package repository

import (
	"fmt"
	"math/rand"
)

func GenerateRandomCodeAndTitle() (string, string) {
	// Generate random alphanumeric strings within the limit of 64 characters
	code := fmt.Sprintf("code_%d", rand.Intn(1000000))   // "code_" + random number
	title := fmt.Sprintf("title_%d", rand.Intn(1000000)) // "title_" + random number
	// Ensure the length is less than or equal to 64 characters
	if len(code) > 64 {
		code = code[:64]
	}
	if len(title) > 64 {
		title = title[:64]
	}

	return code, title
}
func GenerateRandomBool() bool {
	return rand.Intn(2) == 1
}
func GenerateRandomVoucherNumber() string {
	voucher := fmt.Sprintf("voucher_%d", rand.Intn(1000000))

	if len(voucher) > 64 {
		voucher = voucher[:64]
	}
	return voucher
}

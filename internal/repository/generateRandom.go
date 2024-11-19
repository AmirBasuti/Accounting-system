package repository

import (
	"fmt"
	"math/rand"
)

// GenerateRandomCodeAndTitle creates a random code and title
// Both are alphanumeric strings with a maximum length of 64 characters
func GenerateRandomCodeAndTitle() (string, string) {
	// Generate random code and title with "code_" and "title_" prefixes
	code := fmt.Sprintf("code_%d", rand.Intn(1000000))   // Example: code_123456
	title := fmt.Sprintf("title_%d", rand.Intn(1000000)) // Example: title_123456

	// Ensure the length of code and title is not greater than 64 characters
	if len(code) > 64 {
		code = code[:64]
	}
	if len(title) > 64 {
		title = title[:64]
	}

	return code, title
}

// GenerateRandomBool creates a random boolean value (true or false)
func GenerateRandomBool() bool {
	return rand.Intn(2) == 1 // Randomly returns true or false
}

// GenerateRandomVoucherNumber creates a random voucher number
// The result is an alphanumeric string with a maximum length of 64 characters
func GenerateRandomVoucherNumber() string {
	// Generate random voucher number with "voucher_" prefix
	voucher := fmt.Sprintf("voucher_%d", rand.Intn(1000000)) // Example: voucher_123456

	// Ensure the length of the voucher number is not greater than 64 characters
	if len(voucher) > 64 {
		voucher = voucher[:64]
	}

	return voucher
}

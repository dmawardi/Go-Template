package helpers

import (
	"time"

	"github.com/dmawardi/Go-Template/internal/db"
)

// Generates a verification code and sets expiry in user object
func GenerateVerificationCodeAndSetExpiry(user *db.User) error {
	// Generate token for verification
	tokenCode, err := GenerateRandomString(25)
	if err != nil {
		return err
	}
	// Update user to be unverified
	unVerified := false
	user.Verified = &unVerified
	// Set token code
	user.VerificationCode = tokenCode
	// Set verification code expiry to 12 hours from now
	user.VerificationCodeExpiry = time.Now().Add(12 * time.Hour)
	return nil
}

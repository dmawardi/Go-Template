package helpers

import (
	"time"

	"github.com/dmawardi/Go-Template/internal/db"
)

// Generates a verification code and sets expiry and returns a user object with the values
func GenerateVerificationCodeAndSetExpiry() (*db.User, error) {
	userUpdate := &db.User{}
	// Generate token for verification
	tokenCode, err := GenerateRandomString(25)
	if err != nil {
		return nil, err
	}
	// Update user to be unverified
	verified := false
	userUpdate.Verified = &verified
	// Set token code
	userUpdate.VerificationCode = tokenCode
	// Set verification code expiry to 12 hours from now
	userUpdate.VerificationCodeExpiry = time.Now().Add(12 * time.Hour)
	return userUpdate, nil
}

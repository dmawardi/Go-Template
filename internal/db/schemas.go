package db

import (
	"time"

	"gorm.io/gorm"
)

// Schemas
type User struct {
	// gorm.Model `json:"-"`
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `json:"name"`
	Username  string         `json:"username"`
	Email     string         `json:"email" gorm:"uniqueIndex"`
	Password  string         `json:"-"`
	Role      string         `json:"role" gorm:"default:user"`
	// Verification
	Verified               bool      `json:"verified" gorm:"default:false"`
	VerificationCode       string    `json:"verification_code" gorm:"default:null"`
	VerificationCodeExpiry time.Time `json:"verification_code_expiry" gorm:"default:null"`
}

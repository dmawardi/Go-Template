package db

import (
	"time"

	"gorm.io/gorm"
)

// Schemas
type User struct {
	// gorm.Model `json:"-"`
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `swaggertype:"string" json:"created_at,omitempty"`
	UpdatedAt time.Time      `swaggertype:"string" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `json:"name,omitempty"`
	Username  string         `json:"username,omitempty"`
	Email     string         `json:"email,omitempty" gorm:"uniqueIndex"`
	Password  string         `json:"-"`
	Role      string         `json:"role,omitempty" gorm:"default:user"`
	// Verification
	Verified               bool      `json:"verified,omitempty" gorm:"default:false"`
	VerificationCode       string    `json:"verification_code,omitempty" gorm:"default:null"`
	VerificationCodeExpiry time.Time `json:"verification_code_expiry,omitempty" gorm:"default:null"`
	// Relationships
	Posts []Post `json:"posts,omitempty" gorm:"foreignKey:UserID"`
}

type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `swaggertype:"string" json:"created_at,omitempty"`
	UpdatedAt time.Time      `swaggertype:"string" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Title     string         `json:"title,omitempty"`
	Body      string         `json:"body,omitempty"`
	UserID    uint           `json:"user_id,omitempty"`
	User      User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

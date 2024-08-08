package db

import (
	"time"

	"gorm.io/gorm"
)

type Action struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
    AdminID     uint           `gorm:"not null" json:"admin_id"`  // Foreign key to the admin user
    ActionType  string         `gorm:"not null" json:"action_type"`  // Type of action (create, update, delete)
    EntityType  string         `gorm:"not null" json:"entity_type"`  // Type of entity affected (user, product, order, etc.)
    EntityID    uint           `gorm:"not null" json:"entity_id"`  // ID of the affected entity
    Changes     string         `gorm:"type:json" json:"changes"` // JSON field to record the changes made
    Timestamp   time.Time      `gorm:"autoCreateTime" json:"timestamp"`
    IPAddress   string         `gorm:"size:45" json:"ip_address"`   // IP address of the admin
    Description string         `gorm:"type:text" json:"description"` // Description of the action performed
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
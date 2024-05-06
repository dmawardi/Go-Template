package seed

import (
	"fmt"

	"gorm.io/gorm"
)

// Only allows the insertion of unique record into the database
func InsertUniqueRecord(db *gorm.DB, items []interface{}) error {
	for _, item := range items {
		// Attempt to find the existing item or create a new one if not found.
		// This uses all non-zero fields of the item to check for an existing record.
		result := db.FirstOrCreate(item, item)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			fmt.Println("Item already exists, skipping...")
		} else {
			fmt.Println("Seed item created")
		}
	}
	return nil
}

type BasicFactory interface {
	Factory(count int) error
}

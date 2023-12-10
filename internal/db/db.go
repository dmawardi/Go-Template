package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connects to database and returns client
func DbConnect() *gorm.DB {
	// Grab environment variables for connection
	var DB_USER string = os.Getenv("DB_USER")
	var DB_PASS string = os.Getenv("DB_PASS")
	var DB_HOST string = os.Getenv("DB_HOST")
	var DB_PORT string = os.Getenv("DB_PORT")
	var DB_NAME string = os.Getenv("DB_NAME")

	dbUrl := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", DB_HOST, DB_USER, DB_PASS, DB_NAME, DB_PORT)

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&User{}, &Post{})

	return db
}

// Counts the number of records in a table based on conditions
func CountBasedOnConditions(databaseSchema interface{}, conditions []interface{}, dbClient *gorm.DB) (*int64, error) {
	// Fetch metadata from database
	var totalCount int64

	// Count the total number of records
	query := dbClient.Model(databaseSchema)

	// Add conditions to query
	if len(conditions) > 0 {
		// Iterate through conditions (stop at second last element)
		// Increment by 2 to account for condition and value
		for i := 0; i < len(conditions); i += 2 {
			// Extract condition and value
			condition, value := conditions[i].(string), conditions[i+1]
			// For the first condition, use Where
			if i == 0 {
				// Add condition to query
				query = query.Where(condition, value)
			} else {
				// For subsequent conditions, use Or
				query = query.Or(condition, value)
			}
		}
	}

	// Execute query
	countResult := query.Count(&totalCount)
	if countResult.Error != nil {
		return nil, countResult.Error
	}
	return &totalCount, nil
}

// Bulk deletes all records within a table based on ids
func BulkDeleteByIds(databaseSchema interface{}, ids []int, dbClient *gorm.DB) error {
	// Start a transaction (to avoid partial deletion)
	err := dbClient.Transaction(func(tx *gorm.DB) error {
		// In the transaction, delete users with specified IDs
		if err := tx.Where("id IN ?", ids).Delete(&databaseSchema).Error; err != nil {
			return err // Return any error to rollback the transaction
		}

		return nil // Return nil to commit the transaction
	})

	// Check if the transaction was successful
	if err != nil {
		return err
	} else {
		return nil
	}
}

// Query all records in a table based on conditions (results are populated into a slice of the schema type)
func QueryAll(dbClient *gorm.DB, dbSchema interface{}, limit, offset int, order string, conditions []interface{}) error {
	// Build base query for query schema
	query := dbClient.Model(dbSchema)

	// Add parameters into query as needed
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if order != "" {
		// Add order to query
		query = query.Order(order)
	}

	// Iterate through conditions (stop at second last element)
	// Increment by 2 to account for condition and value
	for i := 0; i < len(conditions); i += 2 {
		// Extract condition and value
		condition, value := conditions[i].(string), conditions[i+1]
		// For the first condition, use Where
		if i == 0 {
			// Add condition to query
			query = query.Where(condition, value)
		} else {
			// For subsequent conditions, use Or
			query = query.Or(condition, value)
		}
	}

	// Query database
	if err := query.Find(dbSchema).Error; err != nil {
		fmt.Printf("Error in query all: %v\n", err)
		return err
	}

	return nil
}

// Extract pointer value as string using data type (used in ObtainValue)
func PointerToStringWithType(ptr interface{}, dataType string) string {
	switch dataType {
	case "bool":
		if val, ok := ptr.(*bool); ok {
			if val == nil {
				return "nil"
			}
			return fmt.Sprintf("%t", *val)
		}
	case "int":
		if val, ok := ptr.(*int); ok {
			if val == nil {
				return "nil"
			}
			return fmt.Sprintf("%d", *val)
		}
	case "float64":
		if val, ok := ptr.(*float64); ok {
			if val == nil {
				return "nil"
			}
			return fmt.Sprintf("%f", *val)
		}
	case "string":
		if val, ok := ptr.(*string); ok {
			if val == nil {
				return "nil"
			}
			return *val
		}
	}

	return ""
}

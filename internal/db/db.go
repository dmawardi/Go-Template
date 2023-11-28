package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Post{})

	return db
}

// Counts the number of records in a table based on conditions
func CountBasedOnConditions(databaseSchema interface{}, conditions []string, dbClient *gorm.DB) (*int64, error) {
	// Fetch metadata from database
	var totalCount int64

	// Count the total number of records
	query := dbClient.Model(databaseSchema)

	// Add conditions to query
	if len(conditions) > 0 {
		for _, condition := range conditions {
			// Add condition to query
			query.Where(condition)
		}

	}

	// Execute query
	countResult := query.Count(&totalCount)
	if countResult.Error != nil {
		return nil, countResult.Error
	}
	return &totalCount, nil
}

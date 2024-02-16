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

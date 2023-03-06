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

	return db
}

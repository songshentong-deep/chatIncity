package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "postgres://chat:password@localhost:5432/social_app?sslmode=disable"
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	var result int
	db.Raw("SELECT 1").Scan(&result)
	fmt.Printf("Database connection successful! Result: %d\n", result)
}
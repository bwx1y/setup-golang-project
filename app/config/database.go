package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	dbName := os.Getenv("DB_ADDR")
	db, err := gorm.Open(mysql.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Fatalf("[ Application ] Failed to connect to database: %f", err)
	}

	fmt.Println("[ Application ] Connected to database")
	return db
}

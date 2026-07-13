package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
}

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var database *gorm.DB
	var err error
	for attempt := 1; attempt <= 10; attempt++ {
		database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			DB = database
			return
		}

		log.Printf("Database connection attempt %d failed: %v", attempt, err)
		time.Sleep(3 * time.Second)
	}

	log.Fatal("Error connecting to database: ", err)
}

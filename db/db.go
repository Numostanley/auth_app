package db

import (
	"fmt"
	"log"
	"os"

	"github.com/Numostanley/d8er_app/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBInstance struct {
	DB *gorm.DB
}

var Database DBInstance

func InitDB() {

	godotenv.Load()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Africa/Lagos",
		os.Getenv("PG_HOST"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_DATABASE"),
		os.Getenv("PG_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Postgres failed to setup or connect.", err)
	}

	log.Println("Running Migrations")
	db.AutoMigrate(&models.User{}, &models.Client{})
	log.Println("Migrations Complete")

	Database = DBInstance{
		DB: db,
	}
}

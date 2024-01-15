package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Numostanley/d8er_app/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type CustomLogger struct{}

func (l *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *CustomLogger) Info(context.Context, string, ...interface{}) {}

func (l *CustomLogger) Warn(context.Context, string, ...interface{}) {}

func (l *CustomLogger) Error(context.Context, string, ...interface{}) {}

func (l *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if err != nil && err != gorm.ErrRecordNotFound {
		// Log only if it's an error other than RecordNotFound
		sql, rows := fc()
		log.Printf("%s [%.3fms] %s rows:%d\n", err, float64(time.Since(begin))/float64(time.Millisecond), sql, rows)
	}
}

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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: &CustomLogger{},
	})

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

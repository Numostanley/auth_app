package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Numostanley/d8er_app/env"
	"github.com/Numostanley/d8er_app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type CustomLogger struct{}

func (l *CustomLogger) LogMode(_ logger.LogLevel) logger.Interface {
	return l
}

func (l *CustomLogger) Info(context.Context, string, ...interface{}) {}

func (l *CustomLogger) Warn(context.Context, string, ...interface{}) {}

func (l *CustomLogger) Error(context.Context, string, ...interface{}) {}

func (l *CustomLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Log only if it's an error other than RecordNotFound
		sql, rows := fc()
		log.Printf("%s [%.3fms] %s rows:%d\n", err, float64(time.Since(begin))/float64(time.Millisecond), sql, rows)
	}
}

type Instance struct {
	DB *gorm.DB
}

var Database Instance

func InitDB() {
	enV := env.GetEnv{}
	enV.LoadEnv()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Africa/Lagos",
		enV.PostgresHost,
		enV.PostgresUser,
		enV.PostgresPassword,
		enV.PostgresDB,
		enV.PostgresPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: &CustomLogger{},
	})

	if err != nil {
		log.Fatal("Postgres failed to setup or connect.", err)
	}

	log.Println("Running Migrations")
	if err := db.AutoMigrate(&models.User{}, &models.Client{}, &models.VerificationCode{}); err != nil {
		panic(err)
	}
	log.Println("Migrations Complete")

	Database = Instance{
		DB: db,
	}
}

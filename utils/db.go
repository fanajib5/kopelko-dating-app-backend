package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes and returns a database connection
func InitDB() (*gorm.DB, error) {
	// Define PostgreSQL connection string using environment variables for security
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	timezone := os.Getenv("DB_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		host, user, password, dbname, port, timezone)

	// Configure GORM to use a custom logger for debugging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Log SQL statements that take longer than 1 second
			LogLevel:                  logger.Silent, // Log level: Silent, Error, Warn, Info
			IgnoreRecordNotFoundError: true,          // Ignore "record not found" errors
			Colorful:                  true,
		},
	)

	// Initialize the GORM DB connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set up a connection pool to handle concurrent database connections
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}

	// Configure connection pool settings
	sqlDB.SetMaxOpenConns(20)                  // Max open connections
	sqlDB.SetMaxIdleConns(10)                  // Max idle connections
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Reuse connections for 30 minutes

	log.Println("Database connected successfully")
	return db, nil
}

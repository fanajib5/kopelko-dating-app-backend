package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	loadEnv()

	// Initialize the database connection
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Execute the migration file
	if err := executeSQLFile(db, "databases/seeders/seeder.sql"); err != nil {
		log.Fatalf("Failed to execute migration: %v", err)
	}

	log.Println("Migration completed successfully!")
}

// Function to load environment variables from .env file
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

// Function to initialize the database connection
func initDB() (*sql.DB, error) {
	// Define PostgreSQL connection string using environment variables for security
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbTimezone := os.Getenv("DB_TIMEZONE")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbTimezone)

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

// Function to execute SQL from file
func executeSQLFile(db *sql.DB, filepath string) error {
	// Read SQL file content
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", filepath, err)
	}
	sqlStatements := string(sqlBytes)

	// Execute the SQL statements
	_, err = db.Exec(sqlStatements)
	if err != nil {
		return fmt.Errorf("failed to execute SQL statements: %v", err)
	}
	return nil
}

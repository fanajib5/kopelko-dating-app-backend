package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	APIPort     string
	DatabaseURL string
	JWTSecret   string
	LimitSwipe  int
}

const defaultLimitSwipe = 10

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found, reading from environment variables")
	}

	appEnv := getEnv("APP_ENV", "development")
	apiPort := getEnv("API_PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "")
	if dbURL == "" {
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "postgres")
		pass := getEnv("DB_PASSWORD", "postgres")
		dbname := getEnv("DB_NAME", "dating_app")
		sslmode := getEnv("DB_SSLMODE", "disable")
		dbURL = "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslmode
	}

	jwtSecret := getEnv("JWT_SECRET", "default_secret_key_change_in_production")

	limitSwipeStr := getEnv("LIMIT_SWIPE", "10")
	limitSwipe, err := strconv.Atoi(limitSwipeStr)
	if err != nil {
		limitSwipe = defaultLimitSwipe
	}

	return &Config{
		AppEnv:      appEnv,
		APIPort:     apiPort,
		DatabaseURL: dbURL,
		JWTSecret:   jwtSecret,
		LimitSwipe:  limitSwipe,
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

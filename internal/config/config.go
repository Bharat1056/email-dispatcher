package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL               string
	DBMaxConns          int32
	DBMinConns          int32
	DBMaxConnLifetime   time.Duration
	DBMaxConnIdleTime   time.Duration
	DBHealthCheckPeriod time.Duration
}

// this function is used to load the configuration from the environment variables
// it loads the .env file if found else it loads the environment variables
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, falling back to environment variables")
	}

	return &Config{
		DBURL:               getEnv("DB_URL", "postgres://postgres:password@localhost:5432/queue?sslmode=disable"),
		DBMaxConns:          int32(getEnvAsInt("DB_MAX_CONNS", 25)),
		DBMinConns:          int32(getEnvAsInt("DB_MIN_CONNS", 5)),
		DBMaxConnLifetime:   getEnvAsDuration("DB_MAX_CONN_LIFETIME", 30*time.Minute),
		DBMaxConnIdleTime:   getEnvAsDuration("DB_MAX_CONN_IDLE_TIME", 15*time.Minute),
		DBHealthCheckPeriod: getEnvAsDuration("DB_HEALTH_CHECK_PERIOD", 1*time.Minute),
	}
}

// this function is used to get the environment variable as string else fallback value is returned
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// this function is used to get the environment variable as int else fallback value is returned
func getEnvAsInt(key string, fallback int) int {
	strValue := getEnv(key, "")
	// here strconv used for converting string to int
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return fallback
}

// this function is used to get the environment variable as time.Duration else fallback value is returned
func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	strValue := getEnv(key, "")
	// here time.ParseDuration is used for converting string to time.Duration
	if value, err := time.ParseDuration(strValue); err == nil {
		return value
	}
	return fallback
}

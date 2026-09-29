package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Port               string
	Host               string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleTimeout        time.Duration
	Env                string // development, staging, production
	BaseURL            string
	CORSAllowedOrigins []string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	URI             string // full connection url (alternative to individual fields)
}

type LoggerConfig struct {
	Level  string // debug,info, warn, error
	Format string //json, text
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			Host:         getEnv("HOST", "0.0.0.0"),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDurationEnv("IDLE_TIMEOUT", 60*time.Second),
			Env:          getEnv("APP_ENV", "development"),
			BaseURL:      getEnv("BASE_URL", "http://localhost:8080"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			DBName:          getEnv("DB_NAME", "prepilot"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			URI:             getEnv("DATABASE_URL", "postgresql://postgres:admin@localhost:5432/gopanel?sslmode=disable"),
		},
		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	return &cfg, nil
}

func getEnv(key, rollback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return rollback
}

func getIntEnv(key string, rollback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return rollback
}

func getDurationEnv(key string, rollback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsedVal, err := time.ParseDuration(value); err == nil {
			return parsedVal
		}
	}
	return rollback
}

func (c *DatabaseConfig) GetDatabaseURI() string {
	if c.URI != "" {
		return c.URI
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.SSLMode,
	)
}

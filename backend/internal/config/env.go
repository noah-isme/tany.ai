package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	defaultAppEnv          = "local"
	defaultMaxOpenConns    = 10
	defaultMaxIdleConns    = 5
	defaultConnMaxLifetime = time.Hour
)

// Config contains runtime configuration loaded from environment variables.
type Config struct {
	AppEnv            string
	PostgresURL       string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
}

// Load reads configuration values from the process environment.
func Load() (Config, error) {
	cfg := Config{
		AppEnv:            getEnv("APP_ENV", defaultAppEnv),
		PostgresURL:       os.Getenv("POSTGRES_URL"),
		DBMaxOpenConns:    defaultMaxOpenConns,
		DBMaxIdleConns:    defaultMaxIdleConns,
		DBConnMaxLifetime: defaultConnMaxLifetime,
	}

	if cfg.PostgresURL == "" {
		return Config{}, fmt.Errorf("POSTGRES_URL is required")
	}

	if v := os.Getenv("DB_MAX_OPEN_CONNS"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid DB_MAX_OPEN_CONNS: %w", err)
		}
		cfg.DBMaxOpenConns = parsed
	}

	if v := os.Getenv("DB_MAX_IDLE_CONNS"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid DB_MAX_IDLE_CONNS: %w", err)
		}
		cfg.DBMaxIdleConns = parsed
	}

	if v := os.Getenv("DB_CONN_MAX_LIFETIME"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid DB_CONN_MAX_LIFETIME: %w", err)
		}
		cfg.DBConnMaxLifetime = d
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

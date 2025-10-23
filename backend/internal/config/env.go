package config

import (
	"errors"
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
	defaultAccessTTLMin    = 15
	defaultRefreshTTLDays  = 7
	defaultRefreshCookie   = "__Host_refresh"
	defaultLoginPerMin     = 5
	defaultLoginBurst      = 10
	minJWTSecretLength     = 32
)

// Config contains runtime configuration loaded from environment variables.
type Config struct {
	AppEnv               string
	PostgresURL          string
	DBMaxOpenConns       int
	DBMaxIdleConns       int
	DBConnMaxLifetime    time.Duration
	JWTSecret            string
	AccessTokenTTL       time.Duration
	RefreshTokenTTL      time.Duration
	RefreshCookieName    string
	LoginRateLimitPerMin int
	LoginRateLimitBurst  int
}

// Load reads configuration values from the process environment.
func Load() (Config, error) {
	cfg := Config{
		AppEnv:               getEnv("APP_ENV", defaultAppEnv),
		PostgresURL:          os.Getenv("POSTGRES_URL"),
		DBMaxOpenConns:       defaultMaxOpenConns,
		DBMaxIdleConns:       defaultMaxIdleConns,
		DBConnMaxLifetime:    defaultConnMaxLifetime,
		AccessTokenTTL:       time.Duration(defaultAccessTTLMin) * time.Minute,
		RefreshTokenTTL:      time.Duration(defaultRefreshTTLDays) * 24 * time.Hour,
		RefreshCookieName:    defaultRefreshCookie,
		LoginRateLimitPerMin: defaultLoginPerMin,
		LoginRateLimitBurst:  defaultLoginBurst,
	}

	if cfg.PostgresURL == "" {
		return Config{}, fmt.Errorf("POSTGRES_URL is required")
	}

	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	if len(cfg.JWTSecret) < minJWTSecretLength {
		return Config{}, errors.New("JWT_SECRET must be at least 32 characters")
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

	if v := os.Getenv("ACCESS_TOKEN_TTL_MIN"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid ACCESS_TOKEN_TTL_MIN: %w", err)
		}
		if parsed <= 0 {
			return Config{}, errors.New("ACCESS_TOKEN_TTL_MIN must be greater than zero")
		}
		cfg.AccessTokenTTL = time.Duration(parsed) * time.Minute
	}

	if v := os.Getenv("REFRESH_TOKEN_TTL_DAY"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid REFRESH_TOKEN_TTL_DAY: %w", err)
		}
		if parsed <= 0 {
			return Config{}, errors.New("REFRESH_TOKEN_TTL_DAY must be greater than zero")
		}
		cfg.RefreshTokenTTL = time.Duration(parsed) * 24 * time.Hour
	}

	if v := os.Getenv("REFRESH_COOKIE_NAME"); v != "" {
		cfg.RefreshCookieName = v
	}

	if v := os.Getenv("LOGIN_RATE_LIMIT_PER_MIN"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid LOGIN_RATE_LIMIT_PER_MIN: %w", err)
		}
		if parsed <= 0 {
			return Config{}, errors.New("LOGIN_RATE_LIMIT_PER_MIN must be greater than zero")
		}
		cfg.LoginRateLimitPerMin = parsed
	}

	if v := os.Getenv("LOGIN_RATE_LIMIT_BURST"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid LOGIN_RATE_LIMIT_BURST: %w", err)
		}
		if parsed <= 0 {
			return Config{}, errors.New("LOGIN_RATE_LIMIT_BURST must be greater than zero")
		}
		cfg.LoginRateLimitBurst = parsed
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

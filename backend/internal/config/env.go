package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultAppEnv                = "local"
	defaultMaxOpenConns          = 10
	defaultMaxIdleConns          = 5
	defaultConnMaxLifetime       = time.Hour
	defaultAccessTTLMin          = 15
	defaultRefreshTTLDays        = 7
	defaultRefreshCookie         = "__Host_refresh"
	defaultLoginPerMin           = 5
	defaultLoginBurst            = 10
	defaultStorageDriver         = string(StorageDriverSupabase)
	defaultUploadMaxMB           = 5
	defaultUploadRatePerMin      = 10
	defaultUploadRateBurst       = 10
	defaultKBCacheTTLSeconds     = 60
	defaultKnowledgeRatePer5Min  = 30
	defaultKnowledgeRateBurst    = 30
	defaultChatRatePer5Min       = 30
	defaultChatRateBurst         = 30
	defaultAIModel               = "gemini-1.5-pro"
	minJWTSecretLength           = 32
	defaultExternalHTTPTimeoutMS = 8000
	defaultExternalRateLimitRPM  = 30
)

var defaultAllowedMIMEs = []string{
	"image/jpeg",
	"image/png",
	"image/webp",
	"image/svg+xml",
}

var defaultExternalAllowlist = []string{"noahis.me", "www.noahis.me", "noahisme.vercel.app"}

// Config contains runtime configuration loaded from environment variables.
type Config struct {
	AppEnv                   string
	PostgresURL              string
	DBMaxOpenConns           int
	DBMaxIdleConns           int
	DBConnMaxLifetime        time.Duration
	JWTSecret                string
	AccessTokenTTL           time.Duration
	RefreshTokenTTL          time.Duration
	RefreshCookieName        string
	LoginRateLimitPerMin     int
	LoginRateLimitBurst      int
	Storage                  StorageConfig
	Upload                   UploadConfig
	UploadRateLimitPerMin    int
	UploadRateLimitBurst     int
	KnowledgeCacheTTL        time.Duration
	KnowledgeRateLimitPerMin int
	KnowledgeRateLimitBurst  int
	ChatRateLimitPerMin      int
	ChatRateLimitBurst       int
	ChatModel                string
	AIProvider               string
	GoogleGenAIKey           string
	LeapcellAPIKey           string
	LeapcellProjectID        string
	LeapcellTableID          string
	External                 ExternalConfig
}

// StorageDriver enumerates supported object storage providers.
type StorageDriver string

const (
	StorageDriverSupabase StorageDriver = "supabase"
	StorageDriverS3       StorageDriver = "s3"
)

// StorageConfig captures configuration for object storage integrations.
type StorageConfig struct {
	Driver   StorageDriver
	Supabase SupabaseConfig
	S3       S3Config
}

// SupabaseConfig stores Supabase storage settings.
type SupabaseConfig struct {
	URL         string
	Bucket      string
	ServiceRole string
	PublicURL   string
}

// S3Config stores AWS S3 compatible settings.
type S3Config struct {
	Region          string
	Bucket          string
	Endpoint        string
	PublicBaseURL   string
	AccessKeyID     string
	SecretAccessKey string
	ForcePathStyle  bool
}

// UploadConfig defines upload validation rules.
type UploadConfig struct {
	MaxBytes    int64
	AllowedMIME []string
	AllowSVG    bool
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
		Storage: StorageConfig{
			Driver: StorageDriver(strings.ToLower(getEnv("STORAGE_DRIVER", defaultStorageDriver))),
		},
		Upload: UploadConfig{
			MaxBytes:    int64(defaultUploadMaxMB) * 1024 * 1024,
			AllowedMIME: append([]string{}, defaultAllowedMIMEs...),
			AllowSVG:    false,
		},
		UploadRateLimitPerMin:    defaultUploadRatePerMin,
		UploadRateLimitBurst:     defaultUploadRateBurst,
		KnowledgeCacheTTL:        time.Duration(defaultKBCacheTTLSeconds) * time.Second,
		KnowledgeRateLimitPerMin: perMinuteFromWindow(defaultKnowledgeRatePer5Min),
		KnowledgeRateLimitBurst:  defaultKnowledgeRateBurst,
		ChatRateLimitPerMin:      perMinuteFromWindow(defaultChatRatePer5Min),
		ChatRateLimitBurst:       defaultChatRateBurst,
		ChatModel:                getEnv("GEMINI_MODEL", defaultAIModel),
		AIProvider:               strings.ToLower(getEnv("AI_PROVIDER", "mock")),
		GoogleGenAIKey:           strings.TrimSpace(os.Getenv("GOOGLE_GENAI_API_KEY")),
		External: ExternalConfig{
			HTTPTimeout:     time.Duration(defaultExternalHTTPTimeoutMS) * time.Millisecond,
			DomainAllowlist: append([]string{}, defaultExternalAllowlist...),
			RateLimitRPM:    defaultExternalRateLimitRPM,
		},
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

	if v := os.Getenv("UPLOAD_MAX_MB"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid UPLOAD_MAX_MB: %w", err)
		}
		if parsed <= 0 {
			return Config{}, errors.New("UPLOAD_MAX_MB must be greater than zero")
		}
		cfg.Upload.MaxBytes = int64(parsed) * 1024 * 1024
	}

	if v := os.Getenv("UPLOAD_ALLOWED_MIME"); v != "" {
		tokens := strings.Split(v, ",")
		allowed := make([]string, 0, len(tokens))
		for _, token := range tokens {
			trimmed := strings.TrimSpace(token)
			if trimmed != "" {
				allowed = append(allowed, strings.ToLower(trimmed))
			}
		}
		if len(allowed) == 0 {
			return Config{}, errors.New("UPLOAD_ALLOWED_MIME must contain at least one mime type")
		}
		cfg.Upload.AllowedMIME = allowed
	}

	if v := os.Getenv("ALLOW_SVG"); v != "" {
		enabled, err := strconv.ParseBool(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid ALLOW_SVG: %w", err)
		}
		cfg.Upload.AllowSVG = enabled
	}

	if v := os.Getenv("UPLOAD_RATE_LIMIT_PER_MIN"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid UPLOAD_RATE_LIMIT_PER_MIN: %w", err)
		}
		if parsed <= 0 {
			return Config{}, errors.New("UPLOAD_RATE_LIMIT_PER_MIN must be greater than zero")
		}
		cfg.UploadRateLimitPerMin = parsed
	}

	if v := os.Getenv("UPLOAD_RATE_LIMIT_BURST"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid UPLOAD_RATE_LIMIT_BURST: %w", err)
		}
		if parsed <= 0 {
			return Config{}, errors.New("UPLOAD_RATE_LIMIT_BURST must be greater than zero")
		}
		cfg.UploadRateLimitBurst = parsed
	}

	if v := os.Getenv("KB_CACHE_TTL_SECONDS"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid KB_CACHE_TTL_SECONDS: %w", err)
		}
		if parsed <= 0 {
			return Config{}, errors.New("KB_CACHE_TTL_SECONDS must be greater than zero")
		}
		cfg.KnowledgeCacheTTL = time.Duration(parsed) * time.Second
	}

	if v := os.Getenv("KB_RATE_LIMIT_PER_5MIN"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid KB_RATE_LIMIT_PER_5MIN: %w", err)
		}
		if parsed <= 0 {
			return Config{}, errors.New("KB_RATE_LIMIT_PER_5MIN must be greater than zero")
		}
		cfg.KnowledgeRateLimitPerMin = perMinuteFromWindow(parsed)
		cfg.KnowledgeRateLimitBurst = parsed
	}

	if v := os.Getenv("KB_RATE_LIMIT_BURST"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid KB_RATE_LIMIT_BURST: %w", err)
		}
		if parsed <= 0 {
			return Config{}, errors.New("KB_RATE_LIMIT_BURST must be greater than zero")
		}
		cfg.KnowledgeRateLimitBurst = parsed
	}

	if v := os.Getenv("CHAT_RATE_LIMIT_PER_5MIN"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid CHAT_RATE_LIMIT_PER_5MIN: %w", err)
		}
		if parsed <= 0 {
			return Config{}, errors.New("CHAT_RATE_LIMIT_PER_5MIN must be greater than zero")
		}
		cfg.ChatRateLimitPerMin = perMinuteFromWindow(parsed)
		cfg.ChatRateLimitBurst = parsed
	}

	if v := os.Getenv("CHAT_RATE_LIMIT_BURST"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid CHAT_RATE_LIMIT_BURST: %w", err)
		}
		if parsed <= 0 {
			return Config{}, errors.New("CHAT_RATE_LIMIT_BURST must be greater than zero")
		}
		cfg.ChatRateLimitBurst = parsed
	}

	if v := os.Getenv("AI_MODEL"); v != "" {
		cfg.ChatModel = v
	}

	timeoutMS := defaultExternalHTTPTimeoutMS
	if v := os.Getenv("HTTP_TIMEOUT_MS"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid HTTP_TIMEOUT_MS: %w", err)
		}
		if parsed <= 0 {
			return Config{}, errors.New("HTTP_TIMEOUT_MS must be greater than zero")
		}
		timeoutMS = parsed
	}
	cfg.External.HTTPTimeout = time.Duration(timeoutMS) * time.Millisecond

	if v := os.Getenv("EXTERNAL_RATE_LIMIT_RPM"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid EXTERNAL_RATE_LIMIT_RPM: %w", err)
		}
		if parsed <= 0 {
			return Config{}, errors.New("EXTERNAL_RATE_LIMIT_RPM must be greater than zero")
		}
		cfg.External.RateLimitRPM = parsed
	}

	if v := os.Getenv("EXTERNAL_DOMAIN_ALLOWLIST"); strings.TrimSpace(v) != "" {
		cfg.External.DomainAllowlist = splitAndTrim(v)
	}

	defaults := defaultExternalSources()
	if raw := strings.TrimSpace(os.Getenv("EXTERNAL_SOURCES_DEFAULT")); raw != "" {
		var seeds []ExternalSourceSeed
		if err := json.Unmarshal([]byte(raw), &seeds); err != nil {
			return Config{}, fmt.Errorf("invalid EXTERNAL_SOURCES_DEFAULT: %w", err)
		}
		if len(seeds) > 0 {
			defaults = seeds
		}
	}
	cfg.External.SourcesDefault = defaults

	if err := populateStorageConfig(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func perMinuteFromWindow(per5min int) int {
	if per5min <= 0 {
		return 1
	}
	perMinute := per5min / 5
	if per5min%5 != 0 {
		perMinute++
	}
	if perMinute <= 0 {
		perMinute = 1
	}
	return perMinute
}

func populateStorageConfig(cfg *Config) error {
	switch cfg.Storage.Driver {
	case StorageDriverSupabase:
		supabaseURL := strings.TrimSpace(os.Getenv("SUPABASE_URL"))
		if supabaseURL == "" {
			return errors.New("SUPABASE_URL is required when STORAGE_DRIVER=supabase")
		}
		bucket := strings.TrimSpace(os.Getenv("SUPABASE_BUCKET"))
		if bucket == "" {
			return errors.New("SUPABASE_BUCKET is required when STORAGE_DRIVER=supabase")
		}
		serviceRole := strings.TrimSpace(os.Getenv("SUPABASE_SERVICE_ROLE"))
		if serviceRole == "" {
			return errors.New("SUPABASE_SERVICE_ROLE is required when STORAGE_DRIVER=supabase")
		}

		publicURL := strings.TrimSpace(os.Getenv("SUPABASE_PUBLIC_URL"))
		if publicURL == "" {
			publicURL = strings.TrimSuffix(supabaseURL, "/") + "/storage/v1/object/public/" + bucket
		}

		cfg.Storage.Supabase = SupabaseConfig{
			URL:         strings.TrimSuffix(supabaseURL, "/"),
			Bucket:      bucket,
			ServiceRole: serviceRole,
			PublicURL:   strings.TrimSuffix(publicURL, "/"),
		}
	case StorageDriverS3:
		region := strings.TrimSpace(os.Getenv("S3_REGION"))
		if region == "" {
			return errors.New("S3_REGION is required when STORAGE_DRIVER=s3")
		}
		bucket := strings.TrimSpace(os.Getenv("S3_BUCKET"))
		if bucket == "" {
			return errors.New("S3_BUCKET is required when STORAGE_DRIVER=s3")
		}
		accessKey := strings.TrimSpace(os.Getenv("AWS_ACCESS_KEY_ID"))
		if accessKey == "" {
			return errors.New("AWS_ACCESS_KEY_ID is required when STORAGE_DRIVER=s3")
		}
		secretKey := strings.TrimSpace(os.Getenv("AWS_SECRET_ACCESS_KEY"))
		if secretKey == "" {
			return errors.New("AWS_SECRET_ACCESS_KEY is required when STORAGE_DRIVER=s3")
		}

		endpoint := strings.TrimSpace(os.Getenv("S3_ENDPOINT"))
		publicBase := strings.TrimSpace(os.Getenv("S3_PUBLIC_BASE_URL"))
		forcePathStyle := false
		if v := os.Getenv("S3_FORCE_PATH_STYLE"); v != "" {
			parsed, err := strconv.ParseBool(v)
			if err != nil {
				return fmt.Errorf("invalid S3_FORCE_PATH_STYLE: %w", err)
			}
			forcePathStyle = parsed
		}

		cfg.Storage.S3 = S3Config{
			Region:          region,
			Bucket:          bucket,
			Endpoint:        endpoint,
			PublicBaseURL:   strings.TrimSuffix(publicBase, "/"),
			AccessKeyID:     accessKey,
			SecretAccessKey: secretKey,
			ForcePathStyle:  forcePathStyle,
		}
	default:
		return fmt.Errorf("unsupported STORAGE_DRIVER: %s", cfg.Storage.Driver)
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// ExternalConfig holds configuration for ingesting external knowledge sources.
type ExternalConfig struct {
	SourcesDefault  []ExternalSourceSeed
	HTTPTimeout     time.Duration
	DomainAllowlist []string
	RateLimitRPM    int
}

// ExternalSourceSeed represents default source definitions from configuration.
type ExternalSourceSeed struct {
	Name       string `json:"name"`
	BaseURL    string `json:"base_url"`
	SourceType string `json:"type"`
	Enabled    bool   `json:"enabled"`
}

func defaultExternalSources() []ExternalSourceSeed {
	return []ExternalSourceSeed{
		{
			Name:       "noahis.me",
			BaseURL:    "https://www.noahis.me",
			SourceType: "auto",
			Enabled:    true,
		},
	}
}

func splitAndTrim(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

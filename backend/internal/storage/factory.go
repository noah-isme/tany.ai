package storage

import (
	"fmt"

	"github.com/tanydotai/tanyai/backend/internal/config"
)

// New returns a concrete ObjectStorage implementation based on configuration.
func New(cfg config.StorageConfig) (ObjectStorage, error) {
	switch cfg.Driver {
	case config.StorageDriverSupabase:
		return NewSupabaseStorage(cfg.Supabase)
	case config.StorageDriverS3:
		return NewS3Storage(cfg.S3)
	default:
		return nil, fmt.Errorf("storage: unsupported driver %q", cfg.Driver)
	}
}

package storage

import (
	"context"
	"errors"
)

// ObjectStorage defines the behaviour for storing binary objects and returning public URLs.
type ObjectStorage interface {
	Put(ctx context.Context, key string, content []byte, contentType string) (string, error)
}

// ErrUnsupportedDriver is returned when the configured storage driver is unknown.
var ErrUnsupportedDriver = errors.New("storage: unsupported driver")

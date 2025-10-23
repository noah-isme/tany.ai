package repos

import "errors"

var (
	// ErrNotFound indicates the requested record does not exist.
	ErrNotFound = errors.New("record not found")
	// ErrInvalidSortField indicates the sort field is not supported.
	ErrInvalidSortField = errors.New("invalid sort field")
	// ErrInvalidSortDirection indicates the sort direction is invalid.
	ErrInvalidSortDirection = errors.New("invalid sort direction")
)

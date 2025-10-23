package dto

// PaginationResponse represents a paginated collection response.
type PaginationResponse[T any] struct {
	Items []T   `json:"items"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

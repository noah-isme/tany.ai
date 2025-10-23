package httpapi

import "github.com/gin-gonic/gin"

// ErrorCode enumerates API error identifiers.
type ErrorCode string

const (
	ErrorCodeValidation      ErrorCode = "VALIDATION_ERROR"
	ErrorCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrorCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrorCodeTooManyRequests ErrorCode = "TOO_MANY_REQUESTS"
	ErrorCodeInternal        ErrorCode = "INTERNAL"
)

// ErrorBody represents the standard error payload envelope.
type ErrorBody struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// ErrorResponse is the response wrapper for errors.
type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

// RespondError sends an error response with the standard envelope.
func RespondError(c *gin.Context, status int, code ErrorCode, message string, details interface{}) {
	c.AbortWithStatusJSON(status, ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// RespondData sends a successful response with data payload.
func RespondData(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{"data": data})
}

// RespondList sends a paginated list response.
func RespondList(c *gin.Context, status int, items interface{}, page, limit int, total int64) {
	c.JSON(status, gin.H{
		"items": items,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

package admin

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/tanydotai/tanyai/backend/internal/httpapi"
	"github.com/tanydotai/tanyai/backend/internal/repos"
)

var registerOnce sync.Once

// ensureValidators registers custom validation rules only once.
func ensureValidators() {
	registerOnce.Do(func() {
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			_ = v.RegisterValidation("currency_code", validateCurrencyCode)
		}
	})
}

func validateCurrencyCode(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	if len(value) != 3 {
		return false
	}
	for _, r := range value {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}

func parseListParams(c *gin.Context) repos.ListParams {
	page := parsePositiveInt(c.Query("page"), 1)
	limit := parsePositiveInt(c.Query("limit"), 20)
	if limit > 100 {
		limit = 100
	}

	return repos.ListParams{
		Page:      page,
		Limit:     limit,
		SortField: strings.ToLower(c.Query("sort")),
		SortDir:   strings.ToLower(c.Query("dir")),
	}
}

func parsePositiveInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func respondValidationError(c *gin.Context, err error) {
	httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "invalid request payload", validationDetails(err))
}

func handleListError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, repos.ErrInvalidSortField) {
		httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "invalid sort field", map[string]string{"sort": "unsupported sort field"})
		return true
	}
	if errors.Is(err, repos.ErrInvalidSortDirection) {
		httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "invalid sort direction", map[string]string{"dir": "must be asc or desc"})
		return true
	}
	return handleRepoError(c, err)
}

func handleRepoError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, repos.ErrNotFound) {
		httpapi.RespondError(c, http.StatusNotFound, httpapi.ErrorCodeNotFound, "resource not found", nil)
		return true
	}
	httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "internal server error", nil)
	return true
}

func validationDetails(err error) interface{} {
	if err == nil {
		return nil
	}

	switch v := err.(type) {
	case validator.ValidationErrors:
		details := make(map[string]string, len(v))
		for _, fieldErr := range v {
			details[strings.ToLower(fieldErr.Field())] = fieldErr.Error()
		}
		return details
	case interface {
		Field() string
		Error() string
	}:
		return map[string]string{strings.ToLower(v.Field()): v.Error()}
	default:
		return map[string]string{"message": err.Error()}
	}
}

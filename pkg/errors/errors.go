package errors

import (
	"errors"
	"net/http"
)

// ErrorCode represents different types of errors
type ErrorCode string

const (
	ErrCodeValidation      ErrorCode = "VALIDATION_ERROR"
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrCodeConflict        ErrorCode = "CONFLICT"
	ErrCodeInternal        ErrorCode = "INTERNAL_ERROR"
	ErrCodeBadRequest      ErrorCode = "BAD_REQUEST"
	ErrCodeTooManyRequests ErrorCode = "TOO_MANY_REQUESTS"
)

// APIError represents a structured API error
type APIError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	HTTPStatus int       `json:"-"`
}

func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new API error
func NewAPIError(code ErrorCode, message string, httpStatus int) *APIError {
	return &APIError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// NewAPIErrorWithDetails creates a new API error with details
func NewAPIErrorWithDetails(code ErrorCode, message, details string, httpStatus int) *APIError {
	return &APIError{
		Code:       code,
		Message:    message,
		Details:    details,
		HTTPStatus: httpStatus,
	}
}

// Predefined errors
var (
	ErrUserNotFound    = NewAPIError(ErrCodeNotFound, "User not found", http.StatusNotFound)
	ErrTaskNotFound    = NewAPIError(ErrCodeNotFound, "Task not found", http.StatusNotFound)
	ErrInvalidInput    = NewAPIError(ErrCodeValidation, "Invalid input", http.StatusBadRequest)
	ErrEmailExists     = NewAPIError(ErrCodeConflict, "Email already exists", http.StatusConflict)
	ErrUnauthorized    = NewAPIError(ErrCodeUnauthorized, "Unauthorized", http.StatusUnauthorized)
	ErrForbidden       = NewAPIError(ErrCodeForbidden, "Forbidden", http.StatusForbidden)
	ErrInternalServer  = NewAPIError(ErrCodeInternal, "Internal server error", http.StatusInternalServerError)
	ErrTooManyRequests = NewAPIError(ErrCodeTooManyRequests, "Too many requests", http.StatusTooManyRequests)
)

// WrapError wraps a standard error into an API error
func WrapError(err error, apiErr *APIError) *APIError {
	if err == nil {
		return nil
	}

	// If it's already an APIError, return it
	var apiError *APIError
	if errors.As(err, &apiError) {
		return apiError
	}

	// Wrap the error
	return &APIError{
		Code:       apiErr.Code,
		Message:    apiErr.Message,
		Details:    err.Error(),
		HTTPStatus: apiErr.HTTPStatus,
	}
}

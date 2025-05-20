package errors

import (
	"fmt"
	"net/http"
)

// AppError represents a custom application error
type AppError struct {
	Code        int         `json:"code"`
	Message     string      `json:"message"`
	Details     interface{} `json:"details,omitempty"`
	InternalErr error       `json:"error,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(message string, details interface{}) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Details: details,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource string, id interface{}) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%s with ID %v not found", resource, id),
	}
}

// NewAuthenticationError creates a new authentication error
func NewAuthenticationError(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

// NewInternalError creates a new internal server error
func NewInternalError(err error) *AppError {
	return &AppError{
		Code:        http.StatusInternalServerError,
		Message:     "Internal server error",
		InternalErr: err,
	}
}

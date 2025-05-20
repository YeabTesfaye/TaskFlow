package middleware

import (
	"api/errors"
	"api/logger"
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details any `json:"details,omitempty"`
}

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a custom response writer to catch panics
		srw := &statusResponseWriter{ResponseWriter: w}
		defer func() {
			if err := recover(); err != nil {
				logger.ErrorLogger.Printf("Panic recovered: %v", err)
				appError := errors.NewInternalError(nil)
				respondWithError(srw, appError)
			}
		}()

		// Call the next handler
		next.ServeHTTP(srw, r)
	})
}

// statusResponseWriter wraps http.ResponseWriter to capture the status code
type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func respondWithError(w http.ResponseWriter, err error) {
	// Convert to AppError if it's not already
	var appErr *errors.AppError
	if e, ok := err.(*errors.AppError); ok {
		appErr = e
	} else {
		appErr = errors.NewInternalError(err)
	}

	// Log error details
	logger.ErrorLogger.Printf(
		"Error occurred: Code=%d, Message=%s, Details=%v",
		appErr.Code,
		appErr.Message,
		appErr.Details,
	)

	// Prepare response
	resp := ErrorResponse{
		Code:    appErr.Code,
		Message: appErr.Message,
		Details: appErr.Details,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)
	json.NewEncoder(w).Encode(resp)
}
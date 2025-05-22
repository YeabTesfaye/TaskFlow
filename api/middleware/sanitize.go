package middleware

import (
	"bytes"
	"encoding/json"
	"html"
	"io"
	"net/http"
	"strings"
)

// SanitizeInput is a middleware that sanitizes request input
func SanitizeInput(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only process POST, PUT, PATCH requests
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			// Skip multipart/form-data (file uploads)
			if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
				next.ServeHTTP(w, r)
				return
			}

			// Read the body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusBadRequest)
				return
			}
			// Close the original body
			r.Body.Close()

			// Check if body is JSON
			if isJSON(body) {
				// Parse JSON
				var data any
				if err := json.Unmarshal(body, &data); err != nil {
					http.Error(w, "Invalid JSON", http.StatusBadRequest)
					return
				}

				// Sanitize the data
				sanitizedData := sanitizeInterface(data)

				// Convert back to JSON
				sanitizedBody, err := json.Marshal(sanitizedData)
				if err != nil {
					http.Error(w, "Error processing request", http.StatusInternalServerError)
					return
				}

				// Create new body
				r.Body = io.NopCloser(bytes.NewBuffer(sanitizedBody))
			} else {
				// For non-JSON data, just sanitize as string
				sanitizedBody := sanitizeString(string(body))
				r.Body = io.NopCloser(strings.NewReader(sanitizedBody))
			}
		}

		// Sanitize URL query parameters
		q := r.URL.Query()
		for k, v := range q {
			sanitizedValues := make([]string, len(v))
			for i, value := range v {
				sanitizedValues[i] = sanitizeString(value)
			}
			q[k] = sanitizedValues
		}
		r.URL.RawQuery = q.Encode()

		next.ServeHTTP(w, r)
	})
}

// sanitizeInterface recursively sanitizes all string values in an interface{}
func sanitizeInterface(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			result[key] = sanitizeInterface(value)
		}
		return result
	case []any:
		result := make([]interface{}, len(v))
		for i, value := range v {
			result[i] = sanitizeInterface(value)
		}
		return result
	case string:
		return sanitizeString(v)
	default:
		return v
	}
}

// sanitizeString sanitizes a single string
func sanitizeString(input string) string {
	// First unescape any existing HTML entities
	unescaped := html.UnescapeString(input)

	// Then perform the escape
	escaped := html.EscapeString(unescaped)

	// Trim spaces
	escaped = strings.TrimSpace(escaped)

	// Remove potential SQL injection characters
	escaped = strings.ReplaceAll(escaped, "'", "''")

	return escaped
}

// isJSON checks if the input is JSON
func isJSON(data []byte) bool {
	return json.Valid(data)
}

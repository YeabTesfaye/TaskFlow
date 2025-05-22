package validation

import (
	"errors"
	"regexp"
	"time"
)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(password) > 72 {
		return errors.New("password is too long (maximum 72 characters)")
	}

	// Check for at least one lowercase letter
	if matched, _ := regexp.MatchString(`[a-z]`, password); !matched {
		return errors.New("password must contain at least one lowercase letter")
	}

	// Check for at least one uppercase letter
	if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
		return errors.New("password must contain at least one uppercase letter")
	}

	// Check for at least one digit
	if matched, _ := regexp.MatchString(`\d`, password); !matched {
		return errors.New("password must contain at least one number")
	}

	// Check for at least one special character
	if matched, _ := regexp.MatchString(`[!@#$%^&*(),.?":{}|<>]`, password); !matched {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func IsValidTimezone(tz string) bool {
	_, err := time.LoadLocation(tz)
	return err == nil
}

func IsValidImageType(contentType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}
	return validTypes[contentType]
}

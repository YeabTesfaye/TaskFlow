package utils

import (
	"api/configs"
	"api/errors"
	"api/logger"
	"api/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// RequestValidation validates common request requirements
type RequestValidation struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to hash password: %v", err)
		return "", errors.NewInternalError(err)
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		logger.InfoLogger.Printf("Password check failed for user")
	}
	return err == nil
}

// ValidateRequest validates the common fields in requests
func (rv *RequestValidation) ValidateRequest() error {
	validationErrors := make(map[string]string)

	if rv.Email == "" {
		validationErrors["email"] = "Email is required"
	}
	if rv.Password == "" {
		validationErrors["password"] = "Password is required"
	}

	if len(validationErrors) > 0 {
		return errors.NewValidationError("Validation failed", validationErrors)
	}
	return nil
}

// CheckEmailExists checks if an email already exists in the database
func CheckEmailExists(email string) (bool, error) {
	userCollection := configs.GetCollection(configs.DB, "users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		logger.ErrorLogger.Printf("Database error checking email existence: %v", err)
		return false, errors.NewInternalError(err)
	}
	return count > 0, nil
}

// ValidateUserCreation validates a new user creation request
func ValidateUserCreation(user *models.User) error {
	// Basic validation using the User model's Validate method
	if err := user.Validate(); err != nil {
		return errors.NewValidationError("User validation failed", err)
	}

	// Check if email already exists
	exists, err := CheckEmailExists(user.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.NewValidationError("Email already exists", map[string]string{
			"email": "This email is already registered",
		})
	}

	return nil
}

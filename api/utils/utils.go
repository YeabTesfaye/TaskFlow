package utils

import (
	"api/configs"
	"api/models"
	"context"
	"errors"

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
	// TODO: Implement password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// CheckPassword compares a hashed password with its possible plaintext equivalent

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// ValidateRequest validates the common fields in requests
func (rv *RequestValidation) ValidateRequest() error {
	if rv.Email == "" {
		return errors.New("email is required!")
	}
	if rv.Password == "" {
		return errors.New("password is required")
	}
	return nil

}

// CheckEmailExists checks if an email already exists in the database
func CheckEmailExists(email string) (bool, error) {
	userCollection := configs.GetCollection(configs.DB, "users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Check if email exists
	count, err := userCollection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}
	return count > 0, err
}

// ValidateUserCreation validates a new user creation request
func ValidateUserCreation(user *models.User) error {
	// Basic validation using the User model's Validate method
	if err := user.Validate(); err != nil {
		return err
	}

	// Check if email already exists
	exists, err := CheckEmailExists(user.Email)
	if err != nil {
		return errors.New("error checking email existence")
	}
	if exists {
		return errors.New("email already exists")
	}
	

	return nil
}


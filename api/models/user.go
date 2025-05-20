package models

import (
	"api/validation"
	"errors"
	"net/mail"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

func (u *User) Validate() error {
	// Valiate name
	u.Name = strings.TrimSpace(u.Name)

	if len(u.Name) < 3 || len(u.Name) > 50 {
		return errors.New("Name must be between 3 and 50 characters")
	}

	// Validate email
	u.Email = strings.TrimSpace(u.Email)
	if u.Email == "" {
		return errors.New("email is required")
	}
	_, err := mail.ParseAddress(u.Email)

	if err != nil {
		return errors.New("invalid email format")
	}

	// validate password using the enhanced password validation 
	if err := validation.ValidatePassword(u.Password); err != nil{
		return err
	}
	return nil
}

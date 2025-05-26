package models

import (
	"api/validation"
	"errors"
	"net/mail"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserPreferences struct {
	Timezone           string `json:"timezone" bson:"timezone"`
	EmailNotifications bool   `json:"email_notifications" bson:"email_notifications"`
	PushNotifications  bool   `json:"push_notifications" bson:"push_notifications"`
	DailyDigest       bool   `json:"daily_digest" bson:"daily_digest"`
}
type User struct {
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name              string             `json:"name,omitempty" bson:"name,omitempty"`
	Email             string             `json:"email,omitempty" bson:"email,omitempty"`
	Password          string             `json:"password,omitempty" bson:"password,omitempty"`
	ProfilePicture    string             `json:"profile_picture,omitempty" bson:"profile_picture,omitempty"`
	Preferences       UserPreferences    `json:"preferences" bson:"preferences"`
	EmailVerified     bool               `json:"email_verified" bson:"email_verified"`
	VerificationToken string             `json:"verification_token,omitempty" bson:"verification_token,omitempty"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}

func (u *User) Validate() error {
	// Validate name
	u.Name = strings.TrimSpace(u.Name)
	if len(u.Name) < 3 || len(u.Name) > 50 {
		return errors.New("name must be between 3 and 50 characters")
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

	// Validate password using the enhanced password validation
	if err := validation.ValidatePassword(u.Password); err != nil {
		return err
	}



	return nil
}

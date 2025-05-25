package models

import "time"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UpdateNameRequest struct {
	Name string `json:"name"`
}

type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ProfilePictureResponse struct {
	URL         string    `json:"url"`
	FileName    string    `json:"file_name"`
	ContentType string    `json:"content_type"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserPreferencesUpdate struct {
	Timezone           string `json:"timezone"`
	EmailNotifications bool   `json:"email_notifications"`
	PushNotifications  bool   `json:"push_notifications"`
	DailyDigest       bool   `json:"daily_digest"`
}
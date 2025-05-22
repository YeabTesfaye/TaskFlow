package controllers

import (
	"api/middleware"
	"api/models"
	"api/utils"
	"api/validation"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateProfilePicture(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	
	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userClaims.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID"})
		return
	}

	// Parse multipart form
	err = r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "File too large"})
		return
	}

	file, handler, err := r.FormFile("profile_picture")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid file"})
		return
	}
	defer file.Close()

	// Validate file type
	if !validation.IsValidImageType(handler.Header.Get("Content-Type")) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid file type. Only images are allowed"})
		return
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := "uploads/profile_pictures"
	os.MkdirAll(uploadsDir, os.ModePerm)

	// Generate unique filename using user ID and timestamp
	timestamp := time.Now().Unix()
	fileExt := filepath.Ext(handler.Filename)
	fileName := fmt.Sprintf("%s_%d%s", userClaims.ID, timestamp, fileExt)
	filePath := filepath.Join(uploadsDir, fileName)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save file"})
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save file"})
		return
	}

	// Create URL for accessing the profile picture
	profilePictureURL := fmt.Sprintf("/api/uploads/profile_pictures/%s", fileName)

	// Update user profile picture info in database
	update := bson.M{
		"$set": bson.M{
			"profile_picture": bson.M{
				"file_path": filePath,
				"url": profilePictureURL,
				"updated_at": time.Now(),
				"file_name": fileName,
				"original_name": handler.Filename,
				"content_type": handler.Header.Get("Content-Type"),
			},
			"updated_at": time.Now(),
		},
	}

	result, err := userCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		update,
	)

	if err != nil || result.ModifiedCount == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update profile picture"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Profile picture updated successfully",
		"url": profilePictureURL,
	})
}

func UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userClaims := r.Context().Value("user").(*middleware.UserClaims)

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userClaims.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID"})
		return
	}

	var preferences models.UserPreferences
	// Remove := to use the existing err variable instead of declaring a new one
	err = json.NewDecoder(r.Body).Decode(&preferences)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"preferences": preferences,
			"updated_at":  time.Now(),
		},
	}

	result, err := userCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		update,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update preferences"})
		return
	}

	if result.ModifiedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "No user found to update"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Preferences updated successfully"})
}
func SendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userClaims := r.Context().Value("user").(*middleware.UserClaims)

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userClaims.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID"})
		return
	}

	// Generate verification token
	token := make([]byte, 32)
	if _, err = rand.Read(token); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to generate verification token"})
		return
	}
	verificationToken := hex.EncodeToString(token)

	// Update user with verification token
	update := bson.M{
		"$set": bson.M{
			"verification_token": verificationToken,
			"updated_at":         time.Now(),
		},
	}

	result, err := userCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID}, // Use the converted ObjectID here
		update,
	)

	if err != nil || result.ModifiedCount == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update verification token"})
		return
	}

	// Send verification email using utils
	verificationLink := fmt.Sprintf("http://localhost:8080/api/users/verify-email?token=%s", verificationToken)
	emailBody := fmt.Sprintf("Click <a href=\"%s\">here</a> to verify your email.", verificationLink)

	if err := utils.SendEmail(userClaims.Email, "Verify Your Email", emailBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to send verification email"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Verification email sent successfully",
		"token":   verificationToken, // include this
	})
}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.URL.Query().Get("token")
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Verification token is required"})
		return
	}

	var user models.User
	err := userCollection.FindOne(context.Background(), bson.M{"verification_token": token}).Decode(&user)
	if err != nil {
		fmt.Printf("Error finding user with token: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid or expired verification token"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"email_verified":     true,
			"verification_token": "",
			"updated_at":         time.Now(),
		},
	}

	result, err := userCollection.UpdateOne(
		context.Background(),
		bson.M{"verification_token": token},
		update,
	)

	if err != nil || result.ModifiedCount == 0 {
		fmt.Printf("Error updating user: %v, ModifiedCount: %d\n", err, result.ModifiedCount)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid or expired verification token"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Email verified successfully"})
}

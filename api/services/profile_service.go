package services

import (
	"api/models"
	"api/repositories"
	"api/validation"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileService struct {
	userRepo *repositories.UserRepository
}

func NewProfileService(userRepo *repositories.UserRepository) *ProfileService {
	return &ProfileService{userRepo: userRepo}
}

func (s *ProfileService) UpdateProfilePicture(ctx context.Context, userID primitive.ObjectID, file multipart.File, fileHeader *multipart.FileHeader) (*models.ProfilePictureResponse, error) {
	// Validate file type
	if !validation.IsValidImageType(fileHeader.Header.Get("Content-Type")) {
		return nil, fmt.Errorf("invalid file type. Only images are allowed")
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := "uploads/profile_pictures"
	os.MkdirAll(uploadsDir, os.ModePerm)

	// Generate unique filename
	timestamp := time.Now().Unix()
	fileExt := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%s_%d%s", userID.Hex(), timestamp, fileExt)
	filePath := filepath.Join(uploadsDir, fileName)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	// Create URL for accessing the profile picture
	profilePictureURL := fmt.Sprintf("/api/uploads/profile_pictures/%s", fileName)

	// Update user profile picture info in database
	update := bson.M{
		"profile_picture": bson.M{
			"file_path":     filePath,
			"url":          profilePictureURL,
			"updated_at":    time.Now(),
			"file_name":     fileName,
			"original_name": fileHeader.Filename,
			"content_type":  fileHeader.Header.Get("Content-Type"),
		},
		"updated_at": time.Now(),
	}

	if err := s.userRepo.UpdateUser(ctx, userID, update); err != nil {
		return nil, fmt.Errorf("failed to update profile picture: %v", err)
	}

	return &models.ProfilePictureResponse{
		URL:         profilePictureURL,
		FileName:    fileName,
		ContentType: fileHeader.Header.Get("Content-Type"),
		UpdatedAt:   time.Now(),
	}, nil
}

func (s *ProfileService) UpdatePreferences(ctx context.Context, userID primitive.ObjectID, preferences models.UserPreferencesUpdate) error {
	update := bson.M{
		"preferences": preferences,
		"updated_at": time.Now(),
	}

	return s.userRepo.UpdateUser(ctx, userID, update)
}
package services

import (
	"api/middleware"
	"api/models"
	"api/repositories"
	"api/utils"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (primitive.ObjectID, error) {
	// Validate user
	if err := utils.ValidateUserCreation(user); err != nil {
		return primitive.NilObjectID, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return primitive.NilObjectID, err
	}

	// Set user fields
	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.repo.Create(ctx, user)
}

func (s *UserService) LoginUser(ctx context.Context, email, password string) (*models.User, string, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Check password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate token
	token, err := middleware.GenerateJWT(user.ID.Hex(), user.Email)
	if err != nil {
		return nil, "", err
	}

	// Clear password before returning
	user.Password = ""

	return user, token, nil
}

func (s *UserService) UpdateUserName(ctx context.Context, userID primitive.ObjectID, name string) error {
	update := bson.M{
		"name":       name,
		"updated_at": time.Now(),
	}
	return s.repo.UpdateUser(ctx, userID, update)
}

func (s *UserService) ChangePassword(ctx context.Context, userID primitive.ObjectID, currentPassword, newPassword string) error {
	// Get user
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify current password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	update := bson.M{
		"password":   hashedPassword,
		"updated_at": time.Now(),
	}
	return s.repo.UpdateUser(ctx, userID, update)
}

func (s *UserService) GetUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, error) {
	return s.repo.FindByID(ctx, userID)
}

func (s *UserService) UpdateUser(ctx context.Context, userID primitive.ObjectID, updateData map[string]interface{}) error {
	return s.repo.UpdateUser(ctx, userID, updateData)
}

func (s *UserService) DeleteUser(ctx context.Context, userID primitive.ObjectID) error {
	return s.repo.DeleteUser(ctx, userID)
}


func (s *UserService) SendVerificationEmail(ctx context.Context, userID primitive.ObjectID) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.EmailVerified {
		return errors.New("email is already verified")
	}

	// Generate verification token
	token := primitive.NewObjectID().Hex()

	// Update user with verification token
	update := bson.M{
		"verification_token": token,
		"updated_at":        time.Now(),
	}
	if err := s.repo.UpdateUser(ctx, userID, update); err != nil {
		return err
	}

	// Construct verification URL
	verificationURL := fmt.Sprintf("%s/api/users/verify-email?token=%s", os.Getenv("APP_URL"), token)

	// Prepare email content
	subject := "Verify Your Email"
	htmlBody := fmt.Sprintf(`
		<h2>Email Verification</h2>
		<p>Please click the link below to verify your email address:</p>
		<a href="%s">Verify Email</a>
		<p>If you didn't request this, please ignore this email.</p>
	`, verificationURL)

	return utils.SendEmail(user.Email, subject, htmlBody)
}

func (s *UserService) VerifyEmail(ctx context.Context, token string) error {
	// Find user by verification token
	user, err := s.repo.FindByVerificationToken(ctx, token)
	if err != nil {
		return errors.New("invalid or expired verification token")
	}

	// Update user as verified
	update := bson.M{
		"email_verified":     true,
		"verification_token": "", // Clear the token
		"updated_at":        time.Now(),
	}

	return s.repo.UpdateUser(ctx, user.ID, update)
}
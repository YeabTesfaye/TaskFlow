package controllers

import (
	"api/middleware"
	"api/models"
	"api/services"
	"api/utils"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := c.service.CreateUser(r.Context(), &user)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendJSON(w, map[string]interface{}{
		"id": userID,
	})
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, token, err := c.service.LoginUser(r.Context(), req.Email, req.Password)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	utils.SendJSON(w, models.LoginResponse{
		Token: token,
		User:  *user,
	})
}

func (c *UserController) UpdateName(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	userID, _ := primitive.ObjectIDFromHex(userClaims.ID)

	var req models.UpdateNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.service.UpdateUserName(r.Context(), userID, req.Name); err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, map[string]string{"message": "Name updated successfully"})
}

func (c *UserController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	userID, _ := primitive.ObjectIDFromHex(userClaims.ID)

	var req models.PasswordChangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.service.ChangePassword(r.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendJSON(w, map[string]string{"message": "Password changed successfully"})
}


func (c *UserController) GetMe(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	userID, err := primitive.ObjectIDFromHex(userClaims.ID)
	if err != nil {
		utils.SendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := c.service.GetUserByID(r.Context(), userID)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, user)
}

func (c *UserController) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	userID, err := primitive.ObjectIDFromHex(userClaims.ID)
	if err != nil {
		utils.SendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.service.UpdateUser(r.Context(), userID, updateData); err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, map[string]string{"message": "User updated successfully"})
}

func (c *UserController) DeleteMe(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	userID, err := primitive.ObjectIDFromHex(userClaims.ID)
	if err != nil {
		utils.SendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := c.service.DeleteUser(r.Context(), userID); err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, map[string]string{"message": "User deleted successfully"})
}


func (c *UserController) SendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	userID, err := primitive.ObjectIDFromHex(userClaims.ID)
	if err != nil {
		utils.SendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := c.service.SendVerificationEmail(r.Context(), userID); err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, map[string]string{"message": "Verification email sent successfully"})
}

func (c *UserController) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.SendError(w, "Verification token is required", http.StatusBadRequest)
		return
	}

	if err := c.service.VerifyEmail(r.Context(), token); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendJSON(w, map[string]string{"message": "Email verified successfully"})
}
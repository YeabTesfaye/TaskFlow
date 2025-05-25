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

type ProfileController struct {
	service *services.ProfileService
}

func NewProfileController(service *services.ProfileService) *ProfileController {
	return &ProfileController{service: service}
}

func (c *ProfileController) UpdateProfilePicture(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	userID, err := primitive.ObjectIDFromHex(userClaims.ID)
	if err != nil {
		utils.SendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Parse multipart form
	if err = r.ParseMultipartForm(10 << 20); err != nil {
		utils.SendError(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("profile_picture")
	if err != nil {
		utils.SendError(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	response, err := c.service.UpdateProfilePicture(r.Context(), userID, file, header)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, response)
}

func (c *ProfileController) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	userID, err := primitive.ObjectIDFromHex(userClaims.ID)
	if err != nil {
		utils.SendError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var preferences models.UserPreferencesUpdate
	if err := json.NewDecoder(r.Body).Decode(&preferences); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.service.UpdatePreferences(r.Context(), userID, preferences); err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, map[string]string{"message": "Preferences updated successfully"})
}
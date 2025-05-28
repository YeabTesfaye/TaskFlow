package controllers

import (
	"api/configs"
	"api/middleware"
	"api/models"
	"api/utils"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var tagCollection = configs.GetCollection(configs.DB, "tags")

// CreateTag handles POST /api/tags
func CreateTag(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("user").(*middleware.UserClaims)

	var tag models.Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tag.UserID = userClaims.ID

	if err := tag.Validate(); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := tagCollection.InsertOne(ctx, tag)
	if err != nil {
		utils.SendError(w, "Failed to create tag", http.StatusInternalServerError)
		return
	}

	tag.ID = result.InsertedID.(primitive.ObjectID)
	utils.SendJSON(w, tag)
}

// GetUserTags handles GET /api/tags
func GetUserTags(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("user").(*middleware.UserClaims)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := tagCollection.Find(ctx, bson.M{"user_id": userClaims.ID})
	if err != nil {
		utils.SendError(w, "Failed to fetch tags", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var tags []models.Tag
	if err := cursor.All(ctx, &tags); err != nil {
		utils.SendError(w, "Failed to decode tags", http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, tags)
}

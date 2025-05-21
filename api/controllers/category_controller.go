package controllers

import (
	"api/configs"
	"api/middleware"
	"api/models"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var categoryCollection = configs.GetCollection(configs.DB, "categories")

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	var category models.Category

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	category.UserID = userClaims.ID
	category.CreatedAt = time.Now()

	if err := category.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	result, err := categoryCollection.InsertOne(context.Background(), category)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create category"})
		return
	}

	category.ID = result.InsertedID.(primitive.ObjectID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userClaims := r.Context().Value("user").(*middleware.UserClaims)

	cursor, err := categoryCollection.Find(context.Background(), bson.M{"user_id": userClaims.ID})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch categories"})
		return
	}

	var categories []models.Category
	if err = cursor.All(context.Background(), &categories); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to decode categories"})
		return
	}

	json.NewEncoder(w).Encode(categories)

}

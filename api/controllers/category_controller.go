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
	"go.mongodb.org/mongo-driver/mongo/options"
)

var categoryCollection = configs.GetCollection(configs.DB, "categories")

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	var category models.Category

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	category.UserID = userClaims.ID
	category.CreatedAt = time.Now()

	if err := category.Validate(); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := categoryCollection.InsertOne(context.Background(), category)
	if err != nil {
		utils.SendError(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	category.ID = result.InsertedID.(primitive.ObjectID)
	utils.SendJSON(w, category)
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	params := utils.GetPaginationFromRequest(r)

	filter := bson.M{"user_id": userClaims.ID}
	results, total, err := utils.ExecutePaginatedQuery(r.Context(), categoryCollection, filter, params)
	if err != nil {
		utils.SendError(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	categories := make([]models.Category, 0, len(results))
	for _, result := range results {
		var category models.Category
		bsonBytes, err := bson.Marshal(result)
		if err != nil {
			utils.SendError(w, "Failed to process categories", http.StatusInternalServerError)
			return
		}
		if err := bson.Unmarshal(bsonBytes, &category); err != nil {
			utils.SendError(w, "Failed to process categories", http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}

	response := map[string]interface{}{
		"categories":  categories,
		"page":        params.Page,
		"limit":       params.Limit,
		"total":       total,
		"total_pages": utils.CalculateTotalPages(total, params.Limit),
	}

	utils.SendJSON(w, response)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	categoryID, err := utils.GetObjectIDFromRequest(r, "id")
	if err != nil {
		utils.SendError(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := category.Validate(); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": categoryID, "user_id": userClaims.ID}
	update := bson.M{"$set": bson.M{
		"name":  category.Name,
		"color": category.Color,
	}}

	result := categoryCollection.FindOneAndUpdate(
		r.Context(),
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if result.Err() != nil {
		utils.SendError(w, "Category not found", http.StatusNotFound)
		return
	}

	var updatedCategory models.Category
	if err := result.Decode(&updatedCategory); err != nil {
		utils.SendError(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, updatedCategory)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	categoryID, err := utils.GetObjectIDFromRequest(r, "id")
	if err != nil {
		utils.SendError(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": categoryID, "user_id": userClaims.ID}
	result, err := categoryCollection.DeleteOne(r.Context(), filter)

	if err != nil {
		utils.SendError(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		utils.SendError(w, "Category not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

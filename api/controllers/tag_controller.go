package controllers

import (
	"api/configs"
	"api/models"
	"api/utils"
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var tagCollection = configs.GetCollection(configs.DB, "tags")



func GetTags(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := tagCollection.Find(ctx, bson.M{})
	if err != nil {
		utils.SendError(w, "Failed to fetch tags", http.StatusInternalServerError)
		return
	}

	var tags []models.Tag
	if err := cursor.All(ctx, &tags); err != nil {
		utils.SendError(w, "Failed to decode tags", http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, tags)
}


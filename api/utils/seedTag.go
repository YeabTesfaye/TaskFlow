package utils

import (
	"api/configs"
	"api/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var tagCollection = configs.GetCollection(configs.DB, "tags")

func SeedTags() error {
	initialTags := []models.Tag{
		{Name: "Work", Color: "#3B82F6"},
		{Name: "Personal", Color: "#10B981"},
		{Name: "Urgent", Color: "#EF4444"},
		{Name: "Learning", Color: "#8B5CF6"},
		{Name: "Health", Color: "#EC4899"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Step 1: Delete all existing tags
	if _, err := tagCollection.DeleteMany(ctx, bson.M{}); err != nil {
		return err
	}

	// Step 2: Insert new tags
	var docs []interface{}
	for _, tag := range initialTags {
		docs = append(docs, tag)
	}

	if _, err := tagCollection.InsertMany(ctx, docs); err != nil {
		return err
	}

	return nil
}

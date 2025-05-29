package controllers

import (
	"api/middleware"
	"api/models"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)



func GetTaskStatistics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := context.Background()

	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	userID := userClaims.ID

	// Get total tasks
	totalTasks, err := taskCollection.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch total tasks"}`, http.StatusInternalServerError)
		return
	}

	// Get completed tasks
	completedTasks, err := taskCollection.CountDocuments(ctx, bson.M{
		"user_id": userID,
		"status":  "Completed",
	})
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch completed tasks"}`, http.StatusInternalServerError)
		return
	}

	// Get overdue tasks
	overdueTasks, err := taskCollection.CountDocuments(ctx, bson.M{
		"user_id":  userID,
		"due_date": bson.M{"$lt": time.Now()},
		"status":   bson.M{"$ne": "Completed"},
	})
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch overdue tasks"}`, http.StatusInternalServerError)
		return
	}

	// Calculate completion rate
	var completionRate float64
	if totalTasks > 0 {
		completionRate = float64(completedTasks) / float64(totalTasks) * 100
	}

	// Get task count by priority
	priorityCounts := make(map[string]int)
	cursor, err := taskCollection.Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"user_id": userID}}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":   "$priority",
			"count": bson.M{"$sum": 1},
		}}},
	})
	if err == nil {
		var results []bson.M
		if err = cursor.All(ctx, &results); err == nil {
			for _, result := range results {
				priority, _ := result["_id"].(string)
				count, _ := result["count"].(int32)
				priorityCounts[priority] = int(count)
			}
		}
	}

	// TODO: Add tag-based stats here if needed in future

	stats := models.TaskStatistics{
		ID:             primitive.NewObjectID(),
		UserID:         userID,
		TotalTasks:     int(totalTasks),
		CompletedTasks: int(completedTasks),
		PendingTasks:   int(totalTasks - completedTasks),
		OverdueTasks:   int(overdueTasks),
		CompletionRate: completionRate,
		ByPriority:     priorityCounts,
		UpdatedAt:      time.Now(),
	}

	json.NewEncoder(w).Encode(stats)
}




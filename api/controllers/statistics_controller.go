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
				priority, ok := result["_id"].(string)
				if !ok {
					continue
				}
				count, _ := result["count"].(int32)
				priorityCounts[priority] = int(count)
			}
		}
	}

	// Build and return the statistics
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

func GetTaskStatisticsByCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	ctx := context.Background()

	now := time.Now()

	pipeline := mongo.Pipeline{
		// Match tasks owned by the user
		bson.D{{Key: "$match", Value: bson.M{"user_id": userClaims.ID}}},

		// Unwind categories array to group by each category
		bson.D{{Key: "$unwind", Value: "$categories"}},

		// Group by category ID
		bson.D{{Key: "$group", Value: bson.M{
			"_id":   "$categories",
			"total": bson.M{"$sum": 1},
			"completed": bson.M{"$sum": bson.M{
				"$cond": bson.A{
					bson.M{"$eq": bson.A{"$status", "Completed"}},
					1,
					0,
				},
			}},
			"overdue": bson.M{"$sum": bson.M{
				"$cond": bson.A{
					bson.M{"$and": bson.A{
						bson.M{"$lt": bson.A{"$due_date", now}},
						bson.M{"$ne": bson.A{"$status", "Completed"}},
					}},
					1,
					0,
				},
			}},
		}}},

		// Lookup category details
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "categories",
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "category_info",
		}}},

		// Filter out categories with no match
		bson.D{{Key: "$match", Value: bson.M{
			"category_info": bson.M{"$ne": bson.A{}},
		}}},

		// Project the final result
		bson.D{{Key: "$project", Value: bson.M{
			"category_id": "$_id",
			"category_name": bson.M{
				"$arrayElemAt": bson.A{"$category_info.name", 0},
			},
			"total_tasks":     "$total",
			"completed_tasks": "$completed",
			"completion_rate": bson.M{
				"$cond": bson.M{
					"if":   bson.M{"$eq": bson.A{"$total", 0}},
					"then": 0,
					"else": bson.M{"$multiply": bson.A{
						bson.M{"$divide": bson.A{"$completed", "$total"}},
						100,
					}},
				},
			},
			"overdue_tasks": "$overdue",
		}}},
	}

	cursor, err := taskCollection.Aggregate(ctx, pipeline)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch category statistics"})
		return
	}

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to decode category statistics"})
		return
	}

	json.NewEncoder(w).Encode(results)
}

package controllers

import (
	"api/configs"
	"api/middleware"
	"api/models"
	"context"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskSearchQuery struct {
	Search    string    `json:"search"`     // Search in title and description
	StartDate time.Time `json:"start_date"` // Filter by date range
	EndDate   time.Time `json:"end_date"`
	Priority  string    `json:"priority"` // Filter by priority
	Status    string    `json:"status"`   // Filter by status
	Page      int       `json:"page"`     // Pagination
	Limit     int       `json:"limit"`    // Items per page
}

var taskCollection = configs.GetCollection(configs.DB, "tasks")

func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userClaims := r.Context().Value("user").(*middleware.UserClaims)

	var task models.Task

	_ = json.NewDecoder(r.Body).Decode(&task)

	// Set user ID and timestamps
	task.UserID = userClaims.ID
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	// Validate task input
	if err := task.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := taskCollection.InsertOne(ctx, task)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(result)
}

func GetUserTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	useClaims := r.Context().Value("user").(*middleware.UserClaims)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	cursor, err := taskCollection.Find(ctx, bson.M{"user_id": useClaims.ID})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// GET task ID from URL parameters
	params := mux.Vars(r)
	taskID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid task ID"})
		return
	}

	// Get user claims from context
	userClaims := r.Context().Value("user").(*middleware.UserClaims)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find task by ID and user ID to ensure users can only access their own tasks
	var task models.Task
	err = taskCollection.FindOne(ctx, bson.M{
		"_id":     taskID,
		"user_id": userClaims.ID,
	}).Decode(&task)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Task not found or unauthorized"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(task)
}
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get task ID from URL parameters
	params := mux.Vars(r)
	taskID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid task ID"})
		return
	}

	// Get user claims from context
	userClaims := r.Context().Value("user").(*middleware.UserClaims)

	// Check if the task exists and belongs to the user
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingTask models.Task
	err = taskCollection.FindOne(ctx, bson.M{
		"_id":     taskID,
		"user_id": userClaims.ID,
	}).Decode(&existingTask)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Task not found or unauthorized"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Decode the request body
	var task models.Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Override the user_id with the one from JWT token
	task.UserID = userClaims.ID

	// Validate the task
	if validationErr := task.Validate(); validationErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": validationErr.Error()})
		return
	}

	// Prepare update data
	update := bson.M{
		"$set": bson.M{
			"title":       task.Title,
			"description": task.Description,
			"due_date":    task.DueDate,
			"priority":    task.Priority,
			"status":      task.Status,
			"updated_at":  time.Now(),
		},
	}

	// Update the task
	result, err := taskCollection.UpdateOne(ctx, bson.M{
		"_id":     taskID,
		"user_id": userClaims.ID,
	}, update)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if result.ModifiedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Task not found or unauthorized"})
		return
	}

	// Return the updated task
	var updatedTask models.Task
	err = taskCollection.FindOne(ctx, bson.M{"_id": taskID}).Decode(&updatedTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error retrieving updated task"})
		return
	}

	json.NewEncoder(w).Encode(updatedTask)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	taskID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid task ID"})
		return
	}

	userClaims := r.Context().Value("user").(*middleware.UserClaims)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ensure user can only delete their own tasks
	filter := bson.M{"_id": taskID, "user_id": userClaims.ID}

	result, err := taskCollection.DeleteOne(ctx, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Task not found or unauthorized"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully"})
}

func GetTaskStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userClaims := r.Context().Value("user").(*middleware.UserClaims)

	ctx := context.Background()
	pipeline := []bson.M{
		{"$match": bson.M{"user_id": userClaims.ID}},
		{"$group": bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": 1},
			"completed": bson.M{"$sum": bson.M{
				"$cond": []any{bson.M{"$eq": []string{"$status", "completed"}}, 1, 0},
			}},
			"pending": bson.M{"$sum": bson.M{
				"$cond": []any{bson.M{"$eq": []string{"$status", "pending"}}, 1, 0},
			}},
			"high_priority": bson.M{"$sum": bson.M{
				"$cond": []any{bson.M{"$eq": []string{"$priority", "high"}}, 1, 0},
			}},
		}},
	}

	cursor, err := taskCollection.Aggregate(ctx, pipeline)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch statistics"})
		return
	}

	var stats []bson.M
	if err = cursor.All(ctx, &stats); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to process statistics"})
		return
	}

	if len(stats) == 0 {
		stats = []bson.M{{"total": 0, "completed": 0, "pending": 0, "high_priority": 0}}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats[0])
}

func SearchTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user claims from context
	userClaims := r.Context().Value("user").(*middleware.UserClaims)

	// Parse query parameters
	query := TaskSearchQuery{
		Search:   r.URL.Query().Get("search"),
		Priority: r.URL.Query().Get("priority"),
		Status:   r.URL.Query().Get("status"),
		Page:     1,  // Default values
		Limit:    10, // Default values
	}

	// Parse pagination parameters
	if page := r.URL.Query().Get("page"); page != "" {
		if pageNum, err := strconv.Atoi(page); err == nil && pageNum > 0 {
			query.Page = pageNum
		}
	}
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if limitNum, err := strconv.Atoi(limit); err == nil && limitNum > 0 {
			query.Limit = limitNum
		}
	}

	// Parse date range if provided
	if startDate := r.URL.Query().Get("start_date"); startDate != "" {
		if parsed, err := time.Parse(time.RFC3339, startDate); err == nil {
			query.StartDate = parsed
		}
	}
	if endDate := r.URL.Query().Get("end_date"); endDate != "" {
		if parsed, err := time.Parse(time.RFC3339, endDate); err == nil {
			query.EndDate = parsed
		}
	}

	// Build the filter
	filter := bson.M{"user_id": userClaims.ID}

	// Add search condition if search term is provided
	if query.Search != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": query.Search, "$options": "i"}},
			{"description": bson.M{"$regex": query.Search, "$options": "i"}},
		}
	}

	// Add priority filter if provided
	if query.Priority != "" {
		filter["priority"] = query.Priority
	}

	// Add status filter if provided
	if query.Status != "" {
		filter["status"] = query.Status
	}

	// Add date range filter if provided
	if !query.StartDate.IsZero() || !query.EndDate.IsZero() {
		dateFilter := bson.M{}
		if !query.StartDate.IsZero() {
			dateFilter["$gte"] = query.StartDate
		}
		if !query.EndDate.IsZero() {
			dateFilter["$lte"] = query.EndDate
		}
		filter["due_date"] = dateFilter
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get total count for pagination
	totalCount, err := taskCollection.CountDocuments(ctx, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error counting tasks"})
		return
	}

	// Calculate skip value for pagination
	skip := (query.Page - 1) * query.Limit

	// Find tasks with pagination
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(query.Limit)).
		SetSort(bson.M{"created_at": -1}) // Sort by creation date, newest first

	cursor, err := taskCollection.Find(ctx, filter, opts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error finding tasks"})
		return
	}

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error decoding tasks"})
		return
	}

	// Return response with pagination info
	response := map[string]interface{}{
		"tasks":       tasks,
		"total":       totalCount,
		"page":        query.Page,
		"limit":       query.Limit,
		"total_pages": math.Ceil(float64(totalCount) / float64(query.Limit)),
	}

	json.NewEncoder(w).Encode(response)
}

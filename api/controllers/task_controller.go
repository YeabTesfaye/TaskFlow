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

	// Respond with only the taskId
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"taskId": result.InsertedID.(primitive.ObjectID).Hex(),
	})
}

func GetUserTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user claims from context
	userClaims := r.Context().Value("user").(*middleware.UserClaims)

	// Parse query parameters
	query := r.URL.Query()
	page, _ := strconv.ParseInt(query.Get("page"), 10, 64)
	limit, _ := strconv.ParseInt(query.Get("limit"), 10, 64)
	search := query.Get("search")
	priority := query.Get("priority")
	status := query.Get("status")

	// Set default pagination values
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10 // Default limit
	}

	// Calculate skip value for pagination
	skip := (page - 1) * limit

	// Build filter
	filter := bson.M{"user_id": userClaims.ID}

	// Add search conditions if provided
	if search != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": search, "$options": "i"}},
			{"description": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	// Add priority filter if provided
	if priority != "" {
		filter["priority"] = priority
	}

	// Add status filter if provided
	if status != "" {
		filter["status"] = status
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

	// Configure find options for pagination and sorting
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}). // Sort by creation date, newest first
		SetSkip(skip).
		SetLimit(limit)

	// Execute query
	cursor, err := taskCollection.Find(ctx, filter, opts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error fetching tasks"})
		return
	}

	// Decode results
	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error decoding tasks"})
		return
	}

	// Calculate total pages
	totalPages := int64(math.Ceil(float64(totalCount) / float64(limit)))

	// Prepare response
	response := map[string]interface{}{
		"tasks":       tasks,
		"total":       totalCount,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	}

	json.NewEncoder(w).Encode(response)
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

	w.WriteHeader(http.StatusNoContent)
}



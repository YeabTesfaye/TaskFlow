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
	"go.mongodb.org/mongo-driver/mongo"
)

var taskCollection = configs.GetCollection(configs.DB, "tasks")

// Common task operations
func getTaskByID(ctx context.Context, taskID primitive.ObjectID, userID string) (*models.Task, error) {
	var task models.Task
	err := taskCollection.FindOne(ctx, bson.M{
		"_id":     taskID,
		"user_id": userID,
	}).Decode(&task)
	return &task, err
}

func validateAndPrepareTask(task *models.Task, userID string) error {
	task.UserID = userID
	task.UpdatedAt = time.Now()
	return task.Validate()
}

// HTTP Handlers
func CreateTask(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("user").(*middleware.UserClaims)

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task.CreatedAt = time.Now()
	if err := validateAndPrepareTask(&task, userClaims.ID); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := taskCollection.InsertOne(ctx, task)
	if err != nil {
		utils.SendError(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, map[string]string{
		"taskId": result.InsertedID.(primitive.ObjectID).Hex(),
	})
}

func GetUserTasks(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	params := utils.GetPaginationFromRequest(r)

	baseFilter := bson.M{"user_id": userClaims.ID}
	dateRange := &utils.DateRange{
		StartDate: r.URL.Query().Get("start_date"),
		EndDate:   r.URL.Query().Get("end_date"),
	}

	filter := utils.BuildSearchFilter(baseFilter, params, dateRange)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	results, total, err := utils.ExecutePaginatedQuery(ctx, taskCollection, filter, params)
	if err != nil {
		utils.SendError(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	tasks := make([]models.Task, 0, len(results))
	for _, result := range results {
		var task models.Task
		bsonBytes, err := bson.Marshal(result)
		if err != nil {
			utils.SendError(w, "Failed to process tasks", http.StatusInternalServerError)
			return
		}
		if err := bson.Unmarshal(bsonBytes, &task); err != nil {
			utils.SendError(w, "Failed to process tasks", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	response := map[string]interface{}{
		"tasks":       tasks,
		"total":       total,
		"page":        params.Page,
		"limit":       params.Limit,
		"total_pages": utils.CalculateTotalPages(total, params.Limit),
	}
	utils.SendJSON(w, response)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := utils.GetObjectIDFromRequest(r, "id")
	if err != nil {
		utils.SendError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	task, err := getTaskByID(ctx, taskID, userClaims.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.SendError(w, "Task not found or unauthorized", http.StatusNotFound)
			return
		}
		utils.SendError(w, "Failed to fetch task", http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := utils.GetObjectIDFromRequest(r, "id")
	if err != nil {
		utils.SendError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verify task exists and belongs to user
	_, err = getTaskByID(ctx, taskID, userClaims.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.SendError(w, "Task not found or unauthorized", http.StatusNotFound)
			return
		}
		utils.SendError(w, "Failed to verify task", http.StatusInternalServerError)
		return
	}

	var task models.Task
	if err = json.NewDecoder(r.Body).Decode(&task); err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err = validateAndPrepareTask(&task, userClaims.ID); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"title":       task.Title,
			"description": task.Description,
			"due_date":    task.DueDate,
			"priority":    task.Priority,
			"status":      task.Status,
			"updated_at":  task.UpdatedAt,
		},
	}

	result, err := taskCollection.UpdateOne(ctx, bson.M{
		"_id":     taskID,
		"user_id": userClaims.ID,
	}, update)

	if err != nil {
		utils.SendError(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		utils.SendError(w, "Task not found or unauthorized", http.StatusNotFound)
		return
	}

	updatedTask, err := getTaskByID(ctx, taskID, userClaims.ID)
	if err != nil {
		utils.SendError(w, "Failed to fetch updated task", http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, updatedTask)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := utils.GetObjectIDFromRequest(r, "id")
	if err != nil {
		utils.SendError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := taskCollection.DeleteOne(ctx, bson.M{
		"_id":     taskID,
		"user_id": userClaims.ID,
	})

	if err != nil {
		utils.SendError(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		utils.SendError(w, "Task not found or unauthorized", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetTasksByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := utils.GetObjectIDFromRequest(r, "categoryId")
	if err != nil {
		utils.SendError(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verify category ownership
	var category models.Category
	err = utils.VerifyOwnership(ctx, categoryCollection, categoryID, userClaims.ID, &category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.SendError(w, "Category not found or unauthorized", http.StatusNotFound)
		} else {
			utils.SendError(w, "Failed to verify category", http.StatusInternalServerError)
		}
		return
	}

	params := utils.GetPaginationFromRequest(r)
	baseFilter := bson.M{
		"user_id":    userClaims.ID,
		"categories": categoryID,
	}

	dateRange := &utils.DateRange{
		StartDate: r.URL.Query().Get("start_date"),
		EndDate:   r.URL.Query().Get("end_date"),
	}

	filter := utils.BuildSearchFilter(baseFilter, params, dateRange)
	results, total, err := utils.ExecutePaginatedQuery(ctx, taskCollection, filter, params)
	if err != nil {
		utils.SendError(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	tasks := make([]models.Task, len(results))
	for i, result := range results {
		tasks[i] = result.(models.Task)
	}

	response := map[string]interface{}{
		"tasks":       tasks,
		"page":        params.Page,
		"limit":       params.Limit,
		"total":       total,
		"total_pages": utils.CalculateTotalPages(total, params.Limit),
		"category":    category.Name,
		"category_id": category.ID,
		"color":       category.Color,
		"statistics":  calculateTaskStats(tasks),
	}

	utils.SendJSON(w, response)
}

func calculateTaskStats(tasks []models.Task) map[string]int64 {
	stats := map[string]int64{
		"total":     int64(len(tasks)),
		"completed": 0,
		"pending":   0,
		"overdue":   0,
	}

	now := time.Now()
	for _, task := range tasks {
		switch task.Status {
		case "Completed":
			stats["completed"]++
		case "Pending":
			stats["pending"]++
			if task.DueDate.Before(now) {
				stats["overdue"]++
			}
		}
	}

	return stats
}

package controllers

import (
	"api/configs"
	"api/middleware"
	"api/models"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var commentCollection = configs.GetCollection(configs.DB, "comments")

func AddComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	params := mux.Vars(r)
	taskID := params["taskId"]

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	comment.TaskID = taskID
	comment.UserID = userClaims.ID
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	if err := comment.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	result, err := commentCollection.InsertOne(context.Background(), comment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add comment"})
		return
	}

	comment.ID = result.InsertedID.(primitive.ObjectID)
	json.NewEncoder(w).Encode(comment)
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	taskID := params["taskId"]

	cursor, err := commentCollection.Find(context.Background(), bson.M{"task_id": taskID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch comments"})
		return
	}

	var comments []models.Comment
	if err = cursor.All(context.Background(), &comments); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to decode comments"})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"comments": comments,
	})
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	params := mux.Vars(r)
	commentID, err := primitive.ObjectIDFromHex(params["commentId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid comment ID"})
		return
	}

	var updateComment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&updateComment); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if err := updateComment.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Verify comment ownership
	var existingComment models.Comment
	err = commentCollection.FindOne(context.Background(), bson.M{
		"_id":     commentID,
		"user_id": userClaims.ID,
	}).Decode(&existingComment)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Comment not found or unauthorized"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"content":    updateComment.Content,
			"updated_at": time.Now(),
		},
	}

	_, err = commentCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": commentID},
		update,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update comment"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Comment updated successfully"})
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userClaims := r.Context().Value("user").(*middleware.UserClaims)
	params := mux.Vars(r)
	commentID, err := primitive.ObjectIDFromHex(params["commentId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid comment ID"})
		return
	}

	// Verify comment ownership
	result, err := commentCollection.DeleteOne(context.Background(), bson.M{
		"_id":     commentID,
		"user_id": userClaims.ID,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete comment"})
		return
	}

	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Comment not found or unauthorized"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

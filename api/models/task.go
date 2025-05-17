package models

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	DueDate     time.Time          `json:"due_date" bson:"due_date"`
	Priority    string             `json:"priority" bson:"priority"` // High, Medium, Low
	Status      string             `json:"status" bson:"status"`     // Pending, In Progress, Completed
	UserID      string             `json:"user_id" bson:"user_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

func (t *Task) Valiate() error {
	// Validate Title
	t.Title = strings.TrimSpace(t.Title)
	if len(t.Title) < 3 || len(t.Title) > 100 {
		return errors.New("title must be between 3 and 100 characters")
	}

	// Validate Description
	t.Description = strings.TrimSpace(t.Description)
	if len(t.Description) > 1000 {
		return errors.New("description cannot exceed 1000 characters")
	}

	// Validate Due Date
	if t.DueDate.Before(time.Now()) {
		return errors.New("due date cannot be in the past")
	}

	// Validate Priority
	t.Priority = strings.TrimSpace(t.Priority)
	switch t.Priority {
	case "High", "Medium", "Low":
		// Valid priority
	default:
		return errors.New("priority must be High, Medium, or Low")
	}

	// Validate Status
	t.Status = strings.TrimSpace(t.Status)
	switch t.Status {
	case "Pending", "In Progress", "Completed":
	// valid status
	default:
		return errors.New("status must be Pending, In progress, or Completed")
	}

	// Validate UserID
	if t.UserID == "" {
		return errors.New("user ID is required")
	}
	return nil
}

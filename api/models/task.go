package models

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReminderSettings struct {
	ReminderTime time.Time `json:"reminder_time" bson:"reminder_time"`
	Sent         bool      `json:"sent" bson:"sent"`
}
type Task struct {
	ID          primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string               `json:"title" bson:"title"`
	Description string               `json:"description" bson:"description"`
	DueDate     time.Time            `json:"due_date" bson:"due_date"`
	Priority    string               `json:"priority" bson:"priority"` // High, Medium, Low
	Status      string               `json:"status" bson:"status"`     // Pending, In Progress, Completed
	UserID      string               `json:"user_id" bson:"user_id"`
	Categories  []primitive.ObjectID `json:"categories" bson:"categories"` // category IDs as ObjectIds
	Reminders   []ReminderSettings   `json:"reminders" bson:"reminders"`
	CreatedAt   time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" bson:"updated_at"`
	CompletedAt *time.Time           `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
}

func (t *Task) Validate() error {
	// Title
	t.Title = strings.TrimSpace(t.Title)
	if len(t.Title) < 3 || len(t.Title) > 100 {
		return errors.New("title must be between 3 and 100 characters")
	}

	// Description
	t.Description = strings.TrimSpace(t.Description)
	if len(t.Description) > 1000 {
		return errors.New("description cannot exceed 1000 characters")
	}

	// Due Date
	if t.DueDate.Before(time.Now()) {
		return errors.New("due date cannot be in the past")
	}

	// Priority
	t.Priority = strings.TrimSpace(t.Priority)
	switch t.Priority {
	case "High", "Medium", "Low":
	default:
		return errors.New("priority must be High, Medium, or Low")
	}

	// Status
	t.Status = strings.TrimSpace(t.Status)
	switch t.Status {
	case "Pending", "In Progress", "Completed":
	default:
		return errors.New("status must be Pending, In Progress, or Completed")
	}

	// UserID
	if t.UserID == "" {
		return errors.New("user ID is required")
	}

	// Categories
	if len(t.Categories) > 5 {
		return errors.New("a task can't have more than 5 categories")
	}

	// Categories
	if len(t.Categories) == 0 {
		return errors.New("at least one category is required")
	}
	if len(t.Categories) > 5 {
		return errors.New("a task can't have more than 5 categories")
	}

	// Reminders
	if len(t.Reminders) == 0 {
		return errors.New("at least one reminder is required")
	}

	// Reminders
	for _, r := range t.Reminders {
		if r.ReminderTime.Before(time.Now()) {
			return errors.New("reminder time cannot be in the past")
		}
	}

	return nil
}

package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskStatistics struct {
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID            string             `json:"user_id" bson:"user_id"`
	TotalTasks        int                `json:"total_tasks" bson:"total_tasks"`
	CompletedTasks    int                `json:"completed_tasks" bson:"completed_tasks"`
	PendingTasks      int                `json:"pending_tasks" bson:"pending_tasks"`
	OverdueTasks      int                `json:"overdue_tasks" bson:"overdue_tasks"`
	CompletionRate    float64            `json:"completion_rate" bson:"completion_rate"`
	AverageCompletion float64            `json:"average_completion" bson:"average_completion"` // in days
	ByPriority        map[string]int     `json:"by_priority" bson:"by_priority"`
	ByCategory        map[string]int     `json:"by_category" bson:"by_category"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}

func (ts *TaskStatistics) Validate() error {
	if ts.UserID == "" {
		return errors.New("user ID is required")
	}
	if ts.TotalTasks < 0 || ts.CompletedTasks < 0 || ts.PendingTasks < 0 || ts.OverdueTasks < 0 {
		return errors.New("task counts cannot be negative")
	}
	if ts.CompletionRate < 0 || ts.CompletionRate > 100 {
		return errors.New("completion rate must be between 0 and 100")
	}
	if ts.AverageCompletion < 0 {
		return errors.New("average completion time cannot be negative")
	}
	return nil
}

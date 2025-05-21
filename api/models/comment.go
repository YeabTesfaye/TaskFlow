package models

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	TaskID    string             `json:"task_id" bson:"task_id"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Content   string             `json:"content" bson:"content"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func (c *Comment) Validate() error {
	c.Content = strings.TrimSpace(c.Content)
	if len(c.Content) < 10 || len(c.Content) > 1000 {
		return errors.New("comment must be between 10 and 1000 characters")
	}
	return nil
}

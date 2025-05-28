package models

import (
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
    ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    Name   string             `json:"name" bson:"name"`
    Color  string             `json:"color" bson:"color"`
    UserID string             `json:"user_id" bson:"user_id"`
}

func (t *Tag) Validate() error {
    t.Name = strings.TrimSpace(t.Name)
    if len(t.Name) < 1 || len(t.Name) > 50 {
        return errors.New("tag name must be between 1 and 50 characters")
    }

    t.Color = strings.TrimSpace(t.Color)
    if len(t.Color) < 1 {
        return errors.New("color is required")
    }

    if t.UserID == "" {
        return errors.New("user ID is required")
    }

    return nil
}
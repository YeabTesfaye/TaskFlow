package models

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Color     string             `json:"color" bson:"color"` // Hex color code
	UserID    string             `json:"user_id" bson:"user_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

func (c *Category) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	if len(c.Name) < 2 || len(c.Name) > 30 {
		return errors.New("category name must be between 2 and 30 characters")
	}
	// Validate color format
	c.Color = strings.TrimSpace(c.Color)
	if !strings.HasPrefix(c.Color, "#") || len(c.Color) != 7 {
		return errors.New("color must be a valid hex color code (e.g., #FF0000)")
	}
	return nil
}

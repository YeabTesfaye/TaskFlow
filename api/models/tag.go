package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tag struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Color string             `json:"color" bson:"color"`
}

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email         string `json:"email,omitempty"`
	Notifications []string `json:"notifications,omitempty"`
}
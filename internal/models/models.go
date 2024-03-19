package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todos struct {
	ID		primitive.ObjectID	`json:"_id,omitempty" bson:"_id,omitempty"`
	Title	string				`json:"title,omitempty" bson:"title"`
	Status	*bool				`json:"status" bson:"status"`
}

type CreateTodosRequest struct {
    Title  string `json:"title,omitempty" bson:"title"`
    Status *bool  `json:"status,omitempty" bson:"status"`
}
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Songbook struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Language    string             `json:"language" bson:"language"`
	Country     Country            `json:"country" bson:"country"`
}

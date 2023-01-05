package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Voice struct {
	Voice    string             `json:"voice" bson:"voice"`
	File     primitive.ObjectID `json:"file" bson:"file"`
	Duration float64            `json:"duration" bson:"duration"`
}

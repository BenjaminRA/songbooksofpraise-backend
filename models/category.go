package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID         primitive.ObjectID   `json:"_id" bson:"_id"`
	Category   string               `json:"category" bson:"category"`
	SongbookID primitive.ObjectID   `json:"songbook_id" bson:"songbook_id"`
	ParentID   primitive.ObjectID   `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
	ChildrenID []primitive.ObjectID `json:"children_id,omitempty" bson:"children_id,omitempty"`
}

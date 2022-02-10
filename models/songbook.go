package models

import (
	"context"

	"github.com/BenjaminRA/himnario-backend/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Songbook struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Language    string             `json:"language" bson:"language"`
	Country     Country            `json:"country" bson:"country"`
	Categories  []Category         `json:"categories,omitempty" bson:"categories,omitempty"`
}

func (n *Songbook) GetAllSongbooks() []Songbook {
	db := mongodb.GetMongoDBConnection()

	cursor, err := db.Collection("Songbooks").Aggregate(context.TODO(), []bson.M{
		{"$lookup": bson.M{
			"from":         "Categories",
			"localField":   "_id",
			"foreignField": "songbook_id",
			"pipeline": []bson.M{
				{
					"$match": bson.M{
						"parent_id": primitive.Null{},
					},
				},
				{
					"$project": bson.M{
						"category": 1,
					},
				},
			},
			"as": "categories",
		}},
	})
	if err != nil {
		panic(err)
	}

	result := []Songbook{}

	for cursor.Next(context.TODO()) {
		elem := Songbook{}
		cursor.Decode(&elem)
		result = append(result, elem)
	}

	return result
}

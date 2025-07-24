package models

import (
	"context"

	"github.com/BenjaminRA/himnario-backend/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

type Author struct {
	Author string `json:"author" bson:"author"`
}

func (n *Author) GetAllAuthors() ([]Author, error) {
	db := mongodb.GetMongoDBConnection()

	cursor, err := db.Collection("Authors").Find(context.TODO(), bson.M{})

	if err != nil {
		return []Author{}, err
	}

	var authors []Author
	err = cursor.All(context.TODO(), &authors)

	if err != nil {
		panic(err)
	}

	return authors, nil
}

func AddAuthor(author string) (Author, error) {
	db := mongodb.GetMongoDBConnection()
	newAuthor := Author{}

	cursor := db.Collection("Authors").FindOne(context.TODO(), bson.M{
		"author": author,
	})

	if cursor.Err() != nil {
		_, err := db.Collection("Authors").InsertOne(context.TODO(), bson.M{
			"author": author,
		})

		if err != nil {
			return Author{}, err
		}

		return Author{
			Author: author,
		}, nil
	}

	cursor.Decode(&author)

	return newAuthor, nil
}

package models

import (
	"context"
	"strings"

	"github.com/BenjaminRA/himnario-backend/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BibleBook struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	Code          string             `json:"code" bson:"code"`
	Language_code string             `json:"language_code" bson:"language_code"`
	Book          string             `json:"book" bson:"book"`
	Testament     string             `json:"testament" bson:"testament"`
}

func (n *BibleBook) GetAllBibleBooks(reader_code string) []BibleBook {
	db := mongodb.GetMongoDBConnection()

	if reader_code == "" {
		reader_code = "EN"
	}

	reader_code = strings.ToUpper(reader_code)

	cursor, err := db.Collection("BibleBooks").Find(context.TODO(), bson.M{
		"language_code": reader_code,
	})
	if err != nil {
		panic(err)
	}

	result := []BibleBook{}

	for cursor.Next(context.TODO()) {
		elem := BibleBook{}

		cursor.Decode(&elem)

		result = append(result, elem)
	}

	return result
}

func (n *BibleBook) GetBibleBookByCode(reader_code string, code string) BibleBook {
	db := mongodb.GetMongoDBConnection()

	if reader_code == "" {
		reader_code = "EN"
	}

	reader_code = strings.ToUpper(reader_code)

	cursor, err := db.Collection("Categories").Find(context.TODO(), bson.M{
		"language_code": reader_code,
		"code":          code,
	})
	if err != nil {
		panic(err)
	}

	result := []BibleBook{}

	for cursor.Next(context.TODO()) {
		elem := BibleBook{}
		cursor.Decode(&elem)
		result = append(result, elem)
	}

	if len(result) == 0 {
		return BibleBook{}
	}

	return result[0]
}

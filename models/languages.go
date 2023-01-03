package models

import (
	"context"
	"strings"

	"github.com/BenjaminRA/himnario-backend/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Language struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Language   string             `json:"language" bson:"language"`
	ReaderCode string             `json:"reader_code" bson:"reader_code"`
	Code       string             `json:"code" bson:"code"`
}

func (n *Language) GetAllLanguages(reader_code string) []Language {
	db := mongodb.GetMongoDBConnection()

	if reader_code == "" {
		reader_code = "EN"
	}

	reader_code = strings.ToUpper(reader_code)

	cursor, err := db.Collection("Languages").Find(context.TODO(), bson.M{
		"reader_code": reader_code,
	})
	if err != nil {
		panic(err)
	}

	result := []Language{}

	for cursor.Next(context.TODO()) {
		elem := Language{}

		cursor.Decode(&elem)

		result = append(result, elem)
	}

	return result
}

func (n *Language) GetLanguageByCode(code string, reader_code string) Language {
	db := mongodb.GetMongoDBConnection()

	if reader_code == "" {
		reader_code = "EN"
	}

	reader_code = strings.ToUpper(reader_code)

	cursor := db.Collection("Languages").FindOne(context.TODO(), bson.M{
		"code":        code,
		"reader_code": reader_code,
	})

	result := Language{}

	err := cursor.Decode(&result)

	if err != nil {
		panic(err)
	}

	return result
}

func (n *Language) CreateLanguage(reader_code string) (Language, error) {
	db := mongodb.GetMongoDBConnection()
	n.ID = primitive.NewObjectID()

	if _, err := db.Collection("Languages").InsertOne(context.TODO(), n); err != nil {
		return Language{}, err
	}

	return new(Language).GetLanguageByCode(n.Code, reader_code), nil
}

func (n *Language) UpdateLanguage(code string) error {
	db := mongodb.GetMongoDBConnection()

	_, err := db.Collection("Languages").UpdateMany(context.TODO(), bson.M{
		"code": code,
	}, bson.M{
		"$set": bson.M{
			"language": n.Language,
			"code":     n.Code,
		},
	})

	if err != nil {
		panic(err)
	}

	return nil
}

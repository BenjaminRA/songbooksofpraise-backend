package models

import (
	"context"

	"github.com/BenjaminRA/himnario-backend/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Song struct {
	ID           primitive.ObjectID   `json:"_id" bson:"_id"`
	SongbookID   primitive.ObjectID   `json:"songbook_id" bson:"songbook_id"`
	Categories   []Category           `json:"categories" bson:"categories"`
	CategoriesID []primitive.ObjectID `json:"categories_id" bson:"categories_id"`
	Title        string               `json:"title" bson:"title"`
	Chords       bool                 `json:"chords" bson:"chords"`
	MusicSheet   primitive.ObjectID   `json:"music_sheet" bson:"music_sheet"` //url
	Music        string               `json:"music" bson:"music"`             //url
	Author       string               `json:"author" bson:"author"`
	Description  string               `json:"description" bson:"description"`
	BibleVerse   string               `json:"bible_verse" bson:"bible_verse"`
	Number       int                  `json:"number" bson:"number"`
	Verses       []Verse              `json:"verses" bson:"verses"`
	Voices       []Voice              `json:"voices" bson:"voices"`
}

func (n *Song) GetAllSongs() []Song {
	db := mongodb.GetMongoDBConnection()

	cursor, err := db.Collection("Songs").Aggregate(context.TODO(), []bson.M{
		{"$lookup": bson.M{
			"from":         "Categories",
			"localField":   "categories_id",
			"foreignField": "_id",
			"as":           "categories",
		}},
	})
	if err != nil {
		panic(err)
	}

	result := []Song{}

	for cursor.Next(context.TODO()) {
		elem := Song{}
		cursor.Decode(&elem)
		result = append(result, elem)
	}

	return result
}

func (n *Song) GetSongByID(id string) []Song {
	db := mongodb.GetMongoDBConnection()
	object_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	cursor, err := db.Collection("Songs").Aggregate(context.TODO(), []bson.M{
		{"$match": bson.M{"_id": object_id}},
		{"$lookup": bson.M{
			"from":         "Categories",
			"localField":   "categoriesid",
			"foreignField": "_id",
			"as":           "categories",
		}},
	})
	if err != nil {
		panic(err)
	}

	result := []Song{}

	for cursor.Next(context.TODO()) {
		elem := Song{}
		cursor.Decode(&elem)
		result = append(result, elem)
	}

	return result
}

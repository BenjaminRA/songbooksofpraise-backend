package models

import (
	"bytes"
	"context"
	"errors"
	"log"
	"time"

	"github.com/BenjaminRA/himnario-backend/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
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
	Text         string               `json:"text" bson:"text"`
	Voices       []Voice              `json:"voices" bson:"voices"`
}

func (n *Song) GetAllSongs() []Song {
	db := mongodb.GetMongoDBConnection()

	cursor, err := db.Collection("Songs").Aggregate(context.TODO(), []bson.M{
		{"$lookup": bson.M{
			"from":         "Categories",
			"localField":   "categories_id",
			"foreignField": "_id",
			"pipeline": []bson.M{
				{"$project": bson.M{"category": 1}},
			},
			"as": "categories",
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

func (n *Song) GetSongByID(id string) Song {
	db := mongodb.GetMongoDBConnection()
	object_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	cursor, err := db.Collection("Songs").Aggregate(context.TODO(), []bson.M{
		{"$match": bson.M{"_id": object_id}},
		{"$lookup": bson.M{
			"from":         "Categories",
			"localField":   "categories_id",
			"foreignField": "_id",
			"pipeline": []bson.M{
				{"$project": bson.M{"category": 1, "all": 1}},
			},
			"as": "categories",
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

	return result[0]
}

func (n *Song) GetMusicSheet(id string) ([]byte, string) {
	db := mongodb.GetMongoDBConnection()
	object_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var results bson.M
	err = db.Collection("Songs").FindOne(ctx, bson.M{"_id": object_id}).Decode(&results)
	if err != nil {
		log.Fatal(err)
	}
	// you can print out the result

	var results_object bson.M
	err = db.Collection("fs.files").FindOne(ctx, bson.M{"_id": results["music_sheet"]}).Decode(&results_object)
	if err != nil {
		log.Fatal(err)
	}

	bucket, _ := gridfs.NewBucket(
		db,
	)
	var buf bytes.Buffer
	_, err = bucket.DownloadToStream(results["music_sheet"], &buf)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes(), results_object["filename"].(string)
}

func (n *Song) GetVoice(id string, voice string) ([]byte, string, error) {
	db := mongodb.GetMongoDBConnection()
	object_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var results bson.M
	err = db.Collection("Songs").FindOne(ctx, bson.M{"_id": object_id}).Decode(&results)
	if err != nil {
		log.Fatal(err)
	}

	var file_id primitive.ObjectID

	for _, el := range results["voices"].(bson.A) {
		if el.(bson.M)["voice"].(string) == voice {
			file_id, _ = el.(bson.M)["url"].(primitive.ObjectID)
		}
	}

	if el, _ := primitive.ObjectIDFromHex("000000000000000000000000"); file_id == el {
		return nil, "", errors.New("Voice not found")
	}

	var results_object bson.M
	err = db.Collection("fs.files").FindOne(ctx, bson.M{"_id": file_id}).Decode(&results_object)
	if err != nil {
		log.Fatal(err)
	}

	bucket, _ := gridfs.NewBucket(
		db,
	)
	var buf bytes.Buffer
	_, err = bucket.DownloadToStream(file_id, &buf)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes(), results_object["filename"].(string), nil
}

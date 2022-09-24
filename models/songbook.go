package models

import (
	"context"
	"sync"
	"time"

	"github.com/BenjaminRA/himnario-backend/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Songbook struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Title        string             `json:"title" bson:"title"`
	Description  string             `json:"description" bson:"description"`
	Language     Language           `json:"language,omitempty" bson:"language,omitempty"`
	LanguageCode string             `json:"language_code" bson:"language_code"`
	Country      Country            `json:"country,omitempty" bson:"country,omitempty"`
	CountryCode  string             `json:"country_code" bson:"country_code"`
	Categories   []Category         `json:"categories,omitempty" bson:"categories,omitempty"`
	Numeration   bool               `json:"numeration" bson:"numeration"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

func (n *Songbook) GetAllSongbooks(lang string) []Songbook {
	db := mongodb.GetMongoDBConnection()

	cursor, err := db.Collection("Songbooks").Aggregate(context.TODO(), []bson.M{
		{"$lookup": bson.M{
			"from":         "Languages",
			"localField":   "language_code",
			"foreignField": "code",
			"pipeline": []bson.M{
				{
					"$match": bson.M{
						"reader_code": lang,
					},
				},
			},
			"as": "language",
		}},
		{"$unwind": bson.M{
			"path":                       "$language",
			"preserveNullAndEmptyArrays": true,
		}},
		{"$lookup": bson.M{
			"from":         "Countries",
			"localField":   "country_code",
			"foreignField": "code",
			"pipeline": []bson.M{
				{
					"$match": bson.M{
						"reader_code": lang,
					},
				},
			},
			"as": "country",
		}},
		{"$unwind": bson.M{
			"path":                       "$country",
			"preserveNullAndEmptyArrays": true,
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

func (n *Songbook) GetSongs(id string) []Song {
	db := mongodb.GetMongoDBConnection()
	object_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	cursor, err := db.Collection("Songs").Aggregate(context.TODO(), []bson.M{
		{"$match": bson.M{
			"songbook_id": object_id,
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

func (n *Songbook) GetSongbookByID(id string, lang string) Songbook {
	db := mongodb.GetMongoDBConnection()
	objectID, _ := primitive.ObjectIDFromHex(id)

	cursor, err := db.Collection("Songbooks").Aggregate(context.TODO(), []bson.M{
		{"$match": bson.M{
			"_id": objectID,
		}},
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
					"$sort": bson.M{
						"category": 1,
					},
				},
			},
			"as": "categories",
		}},
		{"$lookup": bson.M{
			"from":         "Languages",
			"localField":   "language_code",
			"foreignField": "code",
			"pipeline": []bson.M{
				{
					"$match": bson.M{
						"reader_code": lang,
					},
				},
			},
			"as": "language",
		}},
		{"$unwind": bson.M{
			"path":                       "$language",
			"preserveNullAndEmptyArrays": true,
		}},
		{"$lookup": bson.M{
			"from":         "Countries",
			"localField":   "country_code",
			"foreignField": "code",
			"pipeline": []bson.M{
				{
					"$match": bson.M{
						"reader_code": lang,
					},
				},
			},
			"as": "country",
		}},
		{"$unwind": bson.M{
			"path":                       "$country",
			"preserveNullAndEmptyArrays": true,
		}},
	})
	if err != nil {
		panic(err)
	}

	result := Songbook{}

	for cursor.Next(context.TODO()) {
		cursor.Decode(&result)

		if result.ID.Hex() != "000000000000000000000000" {
			result.GetCategories()
		}
	}

	return result
}

func (n *Songbook) GetCategories() {
	wg := sync.WaitGroup{}

	for i := 0; i < len(n.Categories); i++ {
		wg.Add(1)

		go func(i int, categories *[]Category) {
			defer wg.Done()
			(*categories)[i].Children = (*categories)[i].GetChildren()
		}(i, &n.Categories)
	}

	AllToFirst(&n.Categories)

	wg.Wait()
}

func (n *Songbook) CreateSongbook(songbook Songbook, lang string) (Songbook, error) {
	db := mongodb.GetMongoDBConnection()

	songbook.ID = primitive.NewObjectID()
	songbook.CreatedAt = time.Now()
	songbook.UpdatedAt = time.Now()

	_, err := db.Collection("Songbooks").InsertOne(context.TODO(), songbook)
	if err != nil {
		return Songbook{}, err
	}

	Category := Category{
		ID:         primitive.NewObjectID(),
		SongbookID: songbook.ID,
		Category:   "Todos",
		All:        true,
	}

	Category.CreateCategory()

	return new(Songbook).GetSongbookByID(songbook.ID.Hex(), lang), nil
}

func (n *Songbook) DeleteSongbook() error {
	db := mongodb.GetMongoDBConnection()

	// Deleting all Categories
	aux := new(Songbook).GetSongbookByID(n.ID.Hex(), "")
	for _, category := range aux.Categories {
		category.DeleteCategory()
	}

	_, err := db.Collection("Songbooks").DeleteOne(context.TODO(), bson.M{
		"_id": n.ID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (n *Songbook) UpdateSongbook() error {
	db := mongodb.GetMongoDBConnection()

	_, err := db.Collection("Songbooks").UpdateOne(context.TODO(), bson.M{
		"_id": n.ID,
	}, bson.M{
		"$set": bson.M{
			"title":         n.Title,
			"description":   n.Description,
			"language_code": n.LanguageCode,
			"country_code":  n.CountryCode,
			"numeration":    n.Numeration,
			"updated_at":    time.Now(),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

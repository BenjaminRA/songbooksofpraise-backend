package migration

import (
	"context"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func addBibleBooks(db *mongo.Database) {
	file, err := os.ReadFile("books.csv")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(file), "\n")

	for i := 1; i < len(lines)-1; i++ {
		row := strings.Split(lines[i], ";")
		db.Collection("BibleBooks").InsertOne(context.TODO(), bson.M{
			"code":          strings.TrimSpace(strings.ToLower(row[0])),
			"language_code": strings.TrimSpace(row[1]),
			"book":          strings.TrimSpace(row[2]),
			"testament":     strings.TrimSpace(row[3]),
		})
	}

}

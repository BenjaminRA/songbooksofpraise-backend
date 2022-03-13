package models

import (
	"context"
	"strings"

	"github.com/BenjaminRA/himnario-backend/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

type Country struct {
	Country    string `json:"country" bson:"country"`
	ReaderCode string `json:"reader_code" bson:"reader_code"`
	Code       string `json:"code" bson:"code"`
}

func (n *Country) GetAllCountries(reader_code string) []Country {
	db := mongodb.GetMongoDBConnection()

	if reader_code == "" {
		reader_code = "EN"
	}

	reader_code = strings.ToUpper(reader_code)

	cursor, err := db.Collection("Countries").Find(context.TODO(), bson.M{
		"reader_code": reader_code,
	})

	if err != nil {
		panic(err)
	}

	var countries []Country
	err = cursor.All(context.TODO(), &countries)

	if err != nil {
		panic(err)
	}

	return countries
}

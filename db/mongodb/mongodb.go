package mongodb

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongodb *mongo.Client

func GetMongoDBConnection() *mongo.Database {
	if mongodb != nil {
		return mongodb.Database("himnario")
	}

	var err error
	err = godotenv.Load()
	if err != nil {
		panic(err)
	}

	credential := options.Credential{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	mongodb, err = mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(
			fmt.Sprintf("mongodb://%s:%s/?readPreference=primary&appname=songbooks_of_praise_backend&directConnection=true&ssl=false",
				os.Getenv("DB_HOST"),
				os.Getenv("DB_PORT"),
			),
		).SetAuth(credential),
	)
	if err != nil {
		panic(err)
	}

	if err := mongodb.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	return mongodb.Database("himnario")
}

func Disconnect() {
	if mongodb != nil {
		if err := mongodb.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
}

func InitDatabase() {
	db := GetMongoDBConnection()

	db.CreateCollection(context.TODO(), "Categories")
	db.CreateCollection(context.TODO(), "Countries")
	db.CreateCollection(context.TODO(), "Languages")
	db.CreateCollection(context.TODO(), "Songs")
	db.CreateCollection(context.TODO(), "Songbooks")
}

func UploadFile(data []byte, filename string) primitive.ObjectID {
	db := GetMongoDBConnection()

	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		panic(err)
	}

	uploadStream, err := bucket.OpenUploadStream(filename)
	if err != nil {
		panic(err)
	}

	defer uploadStream.Close()

	_, err = uploadStream.Write(data)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("Write file to DB was successful. File size: %d M\n", fileSize)

	return uploadStream.FileID.(primitive.ObjectID)
}

func CleanDatabase() {
	db := GetMongoDBConnection()

	db.Drop(context.TODO())
}

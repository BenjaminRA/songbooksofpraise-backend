package mongodb

import (
	"context"

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
	mongodb, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017/?readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false"))
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

	// db.CreateCollection(context.TODO(), "Categories")
	// db.CreateCollection(context.TODO(), "Countries")
	// db.CreateCollection(context.TODO(), "Languages")
	db.CreateCollection(context.TODO(), "Songs")
	db.CreateCollection(context.TODO(), "Songbooks")
	// db.CreateCollection(context.TODO(), "Verses")
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

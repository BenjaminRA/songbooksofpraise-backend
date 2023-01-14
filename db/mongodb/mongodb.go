package mongodb

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongodb *mongo.Client

func GetMongoDBConnection() *mongo.Database {
	if mongodb == nil {
		var err error

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

// Uploads a file to the mongodb database
func UploadFilePath(path string) primitive.ObjectID {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("No se encontró el archivo:", path)
	}

	id := UploadFile(data, fmt.Sprintf("%v.%v", primitive.NewObjectID().Hex(), filepath.Ext(path)))

	return id
}

func FetchFile(collection string, document_id string, field string) ([]byte, string, error) {
	db := GetMongoDBConnection()
	object_id, err := primitive.ObjectIDFromHex(document_id)
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var results bson.M
	err = db.Collection(collection).FindOne(ctx, bson.M{"_id": object_id}).Decode(&results)
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}

	var results_object bson.M
	err = db.Collection("fs.files").FindOne(ctx, bson.M{"_id": results[field]}).Decode(&results_object)
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}

	bucket, _ := gridfs.NewBucket(
		db,
	)
	var buf bytes.Buffer
	_, err = bucket.DownloadToStream(results[field], &buf)
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}
	return buf.Bytes(), results_object["filename"].(string), nil
}

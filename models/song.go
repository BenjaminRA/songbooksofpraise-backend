package models

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/BenjaminRA/himnario-backend/db/mongodb"
	"github.com/tcolgate/mp3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type Song struct {
	ID           primitive.ObjectID   `json:"_id" bson:"_id"`
	SongbookID   primitive.ObjectID   `json:"songbook_id" bson:"songbook_id"`
	Categories   []Category           `json:"categories,omitempty" bson:"categories,omitempty"`
	CategoriesID []primitive.ObjectID `json:"categories_id" bson:"categories_id"`
	Title        string               `json:"title" bson:"title"`
	Chords       bool                 `json:"chords" bson:"chords"`
	MusicSheet   primitive.ObjectID   `json:"music_sheet,omitempty" bson:"music_sheet,omitempty"`
	Music        primitive.ObjectID   `json:"music,omitempty" bson:"music,omitempty"`
	MusicOnly    primitive.ObjectID   `json:"music_only,omitempty" bson:"music_only,omitempty"`
	Author       string               `json:"author,omitempty" bson:"author,omitempty"`
	YouTubeLink  string               `json:"youtube_link,omitempty" bson:"youtube_link,omitempty"`
	Description  string               `json:"description,omitempty" bson:"description,omitempty"`
	BibleVerse   string               `json:"bible_verse,omitempty" bson:"bible_verse,omitempty"`
	Number       int                  `json:"number,omitempty" bson:"number,omitempty"`
	Text         string               `json:"text" bson:"text"`
	Voices       []Voice              `json:"voices,omitempty" bson:"voices,omitempty"`
	CreatedAt    time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at" bson:"updated_at"`
}

func (n *Song) GetAllSongs(args map[string]interface{}) []Song {
	db := mongodb.GetMongoDBConnection()

	cursor, err := db.Collection("Songs").Aggregate(context.TODO(), []bson.M{
		{"$match": args},
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

	if len(result) == 0 {
		return Song{}
	}

	return result[0]
}

func (n *Song) GetMusicSheet(id string) ([]byte, string, error) {
	return mongodb.FetchFile("Songs", id, "music_sheet")
}

func (n *Song) GetMusic(id string) ([]byte, string, error) {
	return mongodb.FetchFile("Songs", id, "music")
}

func (n *Song) GetMusicOnly(id string) ([]byte, string, error) {
	return mongodb.FetchFile("Songs", id, "music_only")
}

func (n *Song) GetVoice(id string, voice string) ([]byte, string, error) {
	db := mongodb.GetMongoDBConnection()
	object_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var results bson.M
	err = db.Collection("Songs").FindOne(ctx, bson.M{"_id": object_id}).Decode(&results)
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}

	var file_id primitive.ObjectID

	if results["voices"] == nil {
		return nil, "", errors.New("this song does not contain voices")
	}

	for _, el := range results["voices"].(bson.A) {
		if el.(bson.M)["voice"].(string) == voice {
			file_id, _ = el.(bson.M)["file"].(primitive.ObjectID)
		}
	}

	if el, _ := primitive.ObjectIDFromHex("000000000000000000000000"); file_id == el {
		return nil, "", errors.New("Voice not found")
	}

	var results_object bson.M
	err = db.Collection("fs.files").FindOne(ctx, bson.M{"_id": file_id}).Decode(&results_object)
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}

	bucket, _ := gridfs.NewBucket(
		db,
	)
	var buf bytes.Buffer
	_, err = bucket.DownloadToStream(file_id, &buf)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", err
	}

	return buf.Bytes(), results_object["filename"].(string), nil
}

func (n *Song) DeleteFileByID(id primitive.ObjectID, db *mongo.Database, ctx context.Context) error {
	_, err := db.Collection("fs.files").DeleteMany(ctx, bson.M{"_id": id})
	return err
}

func (n *Song) SetField(field string, value interface{}, db *mongo.Database) error {
	_, err := db.Collection("Songs").UpdateOne(context.TODO(), bson.M{
		"_id": n.ID,
	}, bson.M{
		"$set": bson.M{
			field:        value,
			"updated_at": time.Now(),
		},
	})

	return err
}

func (n *Song) UnsetField(field string, db *mongo.Database) error {
	_, err := db.Collection("Songs").UpdateOne(context.TODO(), bson.M{
		"_id": n.ID,
	}, bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
		"$unset": bson.M{
			field: "",
		},
	})

	return err
}

func (n *Song) UpdateMusicSheet(path string) error {
	if path == "__same__" {
		return nil
	}

	db := mongodb.GetMongoDBConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if n.ID.Hex() == "000000000000000000000000" {
		return fmt.Errorf("Song not found")
	}

	// Delete current music sheet
	err := n.DeleteFileByID(n.MusicSheet, db, ctx)
	if err != nil {
		return err
	}

	if path != "" {
		n.MusicSheet = mongodb.UploadFilePath(path)
		os.Remove(path)
		n.SetField("music_sheet", n.MusicSheet, db)
	} else {
		n.UnsetField("music_sheet", db)
	}

	return nil
}

func (n *Song) UpdateMusicAudioOnly(path string) error {
	if path == "__same__" {
		return nil
	}

	db := mongodb.GetMongoDBConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if n.ID.Hex() == "000000000000000000000000" {
		return fmt.Errorf("Song not found")
	}

	// Delete current music
	err := n.DeleteFileByID(n.MusicOnly, db, ctx)
	if err != nil {
		return err
	}

	if path != "" {
		n.MusicOnly = mongodb.UploadFilePath(path)
		os.Remove(path)
		n.SetField("music_only", n.MusicOnly, db)
	} else {
		n.UnsetField("music_only", db)
	}

	return nil
}

func (n *Song) UpdateMusicAudio(path string) error {
	if path == "__same__" {
		return nil
	}

	db := mongodb.GetMongoDBConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if n.ID.Hex() == "000000000000000000000000" {
		return fmt.Errorf("Song not found")
	}

	// Delete current music
	err := n.DeleteFileByID(n.Music, db, ctx)
	if err != nil {
		return err
	}

	if path != "" {
		n.Music = mongodb.UploadFilePath(path)
		os.Remove(path)
		n.SetField("music", n.Music, db)
	} else {
		n.UnsetField("music", db)
	}

	return nil
}

func (n *Song) UpdateVoices(path string, voice string) error {
	if path == "__same__" {
		return nil
	}
	db := mongodb.GetMongoDBConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if n.ID.Hex() == "000000000000000000000000" {
		return fmt.Errorf("Song not found")
	}

	// Get voice to edit
	var voice_pointer *Voice
	var voice_index int
	for idx, item := range n.Voices {
		if item.Voice == voice {
			voice_pointer = &item
			voice_index = idx
			break
		}
	}

	if voice_pointer != nil {
		// Delete current voice
		err := n.DeleteFileByID(voice_pointer.File, db, ctx)
		if err != nil {
			return err
		}

		// voice has been deleted
		if path != "" {
			duration := GetFileDuration(path)
			(*voice_pointer).File = mongodb.UploadFilePath(path)
			(*voice_pointer).Duration = duration
			os.Remove(path)
			n.SetField("voices", n.Voices, db)
			// Voice has been replaced
		} else {
			if len(n.Voices) == 1 {
				n.UnsetField("voices", db)
				n.Voices = nil
			} else {
				copy(n.Voices[voice_index:], n.Voices[voice_index+1:])
				n.Voices = n.Voices[:len(n.Voices)-1]
				n.SetField("voices", n.Voices, db)
			}
		}
		// This voices has just been added
	} else {
		if n.Voices == nil {
			n.Voices = []Voice{}
		}

		duration := GetFileDuration(path)
		new_voice := Voice{
			Voice:    voice,
			File:     mongodb.UploadFilePath(path),
			Duration: duration,
		}
		os.Remove(path)

		n.Voices = append(n.Voices, new_voice)
		n.SetField("voices", n.Voices, db)
	}

	return nil
}

func (n *Song) UpdateSong() error {
	db := mongodb.GetMongoDBConnection()

	_, err := db.Collection("Songs").UpdateOne(context.TODO(), bson.M{
		"_id": n.ID,
	}, bson.M{
		"$set": bson.M{
			"songbook_id":   n.SongbookID,
			"categories_id": n.CategoriesID,
			"title":         n.Title,
			"chords":        n.Chords,
			"author":        n.Author,
			"youtube_link":  n.YouTubeLink,
			"description":   n.Description,
			"bible_verse":   n.BibleVerse,
			"number":        n.Number,
			"text":          n.Text,
			"updated_at":    time.Now(),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (n *Song) CreateSong() error {
	db := mongodb.GetMongoDBConnection()

	n.ID = primitive.NewObjectID()
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()

	_, err := db.Collection("Songs").InsertOne(context.TODO(), n)
	if err != nil {
		return err
	}

	return nil
}

func GetFileDuration(path string) float64 {
	t := 0.0

	r, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	d := mp3.NewDecoder(r)
	var f mp3.Frame
	skipped := 0

	for {
		if err := d.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		t = t + f.Duration().Seconds()
	}

	return t
}

func CleanUpSongTemp(song map[string]interface{}) {
	files := []string{"music_sheet_path",
		"music_audio_path",
		"music_audio_only_path",
		"soprano_voice_audio_path",
		"contralto_voice_audio_path",
		"tenor_voice_audio_path",
		"bass_voice_audio_path",
		"all_voice_audio_path",
	}
	for _, file := range files {
		if song[file] != nil {
			value := song[file].(string)
			os.Remove(value)
		}
	}
}

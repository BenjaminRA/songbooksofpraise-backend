package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/BenjaminRA/himnario-backend/db/mongodb"
	"github.com/BenjaminRA/himnario-backend/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tcolgate/mp3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Helper function to recursively print the Temas.
func printTemas(tema models.Tema, counter ...int) {
	if len(counter) == 0 {
		fmt.Println("Tema:", tema.Tema)
	} else {
		fmt.Println("\t", "SubTema:", tema.Tema)
	}

	if len(tema.SubTemas) > 0 {
		for _, subtema := range tema.SubTemas {
			printTemas(subtema, 2)
		}
	} else {
		for _, himno := range tema.Himnos {
			if len(counter) == 0 {
				fmt.Println("\t", himno.ID, himno.Titulo)
			} else {
				fmt.Println("\t\t", himno.ID, himno.Titulo)
			}
		}
	}
}

// Converts an himno model from the old database to the new Song model.
//
// It extracts the music_sheet and voices of hymns and stores them in the new database with the new model.
func HimnoToSong(himno *models.Himno, songbook_id primitive.ObjectID, categories []primitive.ObjectID) models.Song {
	chords := true
	new_verses := ""
	number := 1
	for _, parrafo := range himno.Parrafos {
		if parrafo.Coro {
			new_verses += "{Chorus}\n"
		} else {
			new_verses += fmt.Sprintf("{Verse %d}\n", number)
			number++
		}
		new_verses += fmt.Sprintf("%s\n\n", parrafo.Parrafo)
	}

	music_sheet := uploadFile(fmt.Sprintf("./assets/hymns/%v.jpg", himno.ID))
	voices := []string{"Bajo", "ContraAlto", "Soprano", "Tenor", "Todos"}
	voices_object := []models.Voice{}

	all_voices := true
	for _, voice := range voices {
		if _, err := os.Stat(fmt.Sprintf("./assets/voices/%v/%v.mp3", himno.ID, voice)); err != nil {
			all_voices = false
			break
		}
	}

	if all_voices {
		for _, voice := range voices {
			path := fmt.Sprintf("./assets/voices/%v/%v.mp3", himno.ID, voice)
			id := uploadFile(path)
			voices_object = append(voices_object, models.Voice{
				Voice:    voice,
				File:     id,
				Duration: getFileDuration(path),
			})
		}
	}

	// If the himno id is greater than 517 in the old database, it means is a Coro
	if himno.ID > 517 {
		himno.ID = himno.ID - 517
	}

	return models.Song{
		ID:           primitive.NewObjectID(),
		Number:       himno.ID,
		Title:        himno.Titulo,
		Chords:       chords,
		Text:         new_verses,
		SongbookID:   songbook_id,
		CategoriesID: categories,
		MusicSheet:   music_sheet,
		Voices:       voices_object,
		Music:        "",
		Author:       "",
		Description:  "",
		BibleVerse:   "",
	}
}

// Checks if whether the element ID exists in the array.
func contains(array []primitive.ObjectID, element primitive.ObjectID) bool {
	for _, el := range array {
		if el == element {
			return true
		}
	}

	return false
}

// Checks if a given song has already been added to the database.
func checkIfExists(db *mongo.Database, song *models.Himno) bool {
	itemCount, err := db.Collection("Songs").CountDocuments(context.TODO(), bson.M{"number": song.ID})
	if err != nil {
		panic(err)
	}

	return itemCount > 0
}

// Calculates the total duration in seconds of an MP3 file.
func getFileDuration(path string) float64 {
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

// Adds a category to a specific himno in the himnos_tema dictionary.
func addCategoryToHimno(himnos_tema *map[int]([]primitive.ObjectID), himno_id int, category_id primitive.ObjectID) {
	// If the himno doesn't exists in the himnos_tema dictionary, we initialize the key with an empty array as a value
	if _, found := (*himnos_tema)[himno_id]; !found {
		(*himnos_tema)[himno_id] = []primitive.ObjectID{}
	}

	// If the array corresponding to the value in the himnos_tema dictionary, check whether the category has been added to the array.
	// If not, add it to the array.
	if !contains((*himnos_tema)[himno_id], category_id) {
		(*himnos_tema)[himno_id] = append((*himnos_tema)[himno_id], category_id)
	}
}

// Uploads a file to the mongodb database
func uploadFile(path string) primitive.ObjectID {
	data, err := ioutil.ReadFile(path)
	id := primitive.ObjectID{}
	if err != nil {
		fmt.Println("No se encontró el archivo:", path)
	}

	id = mongodb.UploadFile(data, fmt.Sprintf("%v.%v", primitive.NewObjectID().Hex(), filepath.Ext(path)))

	return id
}

func Migrate() {
	// Initializng database
	mongodb.CleanDatabase()
	mongodb.InitDatabase()

	db := mongodb.GetMongoDBConnection()

	db.Collection("Languages").InsertOne(context.TODO(), bson.M{
		"code":        "ES",
		"reader_code": "ES",
		"language":    "Español",
	})

	db.Collection("Languages").InsertOne(context.TODO(), bson.M{
		"code":        "ES",
		"reader_code": "EN",
		"language":    "Spanish",
	})

	db.Collection("Languages").InsertOne(context.TODO(), bson.M{
		"code":        "EN",
		"reader_code": "ES",
		"language":    "Ingles",
	})

	db.Collection("Languages").InsertOne(context.TODO(), bson.M{
		"code":        "EN",
		"reader_code": "EN",
		"language":    "English",
	})

	db.Collection("Countries").InsertOne(context.TODO(), bson.M{
		"code":        "CL",
		"reader_code": "EN",
		"country":     "Chile",
	})

	db.Collection("Countries").InsertOne(context.TODO(), bson.M{
		"code":        "CL",
		"reader_code": "ES",
		"country":     "Chile",
	})

	songbook := models.Songbook{
		ID:           primitive.NewObjectID(),
		Title:        "Himnos y Cánticos del Evangelio",
		LanguageCode: "ES",
		Description:  "...",
		CountryCode:  "CL",
		Numeration:   true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	songbook_result, _ := db.Collection("Songbooks").InsertOne(context.TODO(), songbook)

	// Getting all hymns
	himnos, _ := new(models.Himno).GetHimnos()

	// Getting verses of hymns
	for i, himno := range himnos {
		parrafos, _ := new(models.Parrafo).GetParrafos(himno.ID)
		himnos[i].Parrafos = parrafos
	}

	// Get all categories
	temas, _ := new(models.Tema).GetAllTemas()

	// Get hymns categories
	himnos_category := map[int]([]primitive.ObjectID){}

	temas = append(temas, models.Tema{
		InsertedID: primitive.NewObjectID(),
		Tema:       "Todos",
		Himnos:     himnos,
	})

	// Getting all subcategories and hymns of every category and/or subcategory
	for i := 0; i < len(temas); i++ {

		temas[i].InsertedID = primitive.NewObjectID()
		temas[i].Himnos, _ = temas[i].GetHimnos()
		temas[i].SubTemas, _ = temas[i].GetAllSubTemas(temas[i].ID)
		for j := 0; j < len(temas[i].Himnos); j++ {
			addCategoryToHimno(&himnos_category, temas[i].Himnos[j].ID, temas[i].InsertedID)
			temas[i].Himnos[j].Parrafos, _ = new(models.Parrafo).GetParrafos(temas[i].Himnos[j].ID)
		}

		for j := 0; j < len(temas[i].SubTemas); j++ {
			temas[i].SubTemas[j].InsertedID = primitive.NewObjectID()
			temas[i].SubTemas[j].Himnos, _ = temas[i].SubTemas[j].GetSubTemaHimnos()
			for k := 0; k < len(temas[i].SubTemas[j].Himnos); k++ {
				addCategoryToHimno(&himnos_category, temas[i].SubTemas[j].Himnos[k].ID, temas[i].SubTemas[j].InsertedID)
				temas[i].SubTemas[j].Himnos[k].Parrafos, _ = new(models.Parrafo).GetParrafos(temas[i].SubTemas[j].Himnos[k].ID)
			}
		}

	}

	var todos_id primitive.ObjectID
	for _, tema := range temas {
		if tema.Tema == "Todos" {
			todos_id = tema.InsertedID
		}
		db.Collection("Categories").InsertOne(context.TODO(), models.Category{
			ID:         tema.InsertedID,
			All:        (tema.Tema == "Todos"),
			SongbookID: songbook_result.InsertedID.(primitive.ObjectID),
			Category:   tema.Tema,
		})

		// children := []primitive.ObjectID{}
		if len(tema.SubTemas) > 0 {
			for _, subtema := range tema.SubTemas {
				// children = append(children, subtema.InsertedID)
				db.Collection("Categories").InsertOne(context.TODO(), models.Category{
					ID:         subtema.InsertedID,
					All:        (subtema.Tema == "Todos"),
					Category:   subtema.Tema,
					SongbookID: songbook_result.InsertedID.(primitive.ObjectID),
					ParentID:   tema.InsertedID,
				})
			}
		}
	}

	for _, himno := range himnos {
		addCategoryToHimno(&himnos_category, himno.ID, todos_id)
		db.Collection("Songs").InsertOne(context.TODO(), HimnoToSong(
			&himno,
			songbook_result.InsertedID.(primitive.ObjectID),
			himnos_category[himno.ID],
		))
	}

	db.Collection("Songbooks").InsertOne(context.TODO(), songbook)

	songbook = models.Songbook{
		ID:           primitive.NewObjectID(),
		Title:        "Coros",
		Numeration:   false,
		LanguageCode: "ES",
		Description:  "...",
		CountryCode:  "CL",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	songbook_result, _ = db.Collection("Songbooks").InsertOne(context.TODO(), songbook)

	// Getting all coros
	himnos, _ = new(models.Himno).GetCoros()

	tema_coros := models.Tema{
		InsertedID: primitive.NewObjectID(),
		Tema:       "Todos",
		Himnos:     himnos,
	}

	db.Collection("Categories").InsertOne(context.TODO(), models.Category{
		ID:         tema_coros.InsertedID,
		All:        true,
		SongbookID: songbook_result.InsertedID.(primitive.ObjectID),
		Category:   tema_coros.Tema,
	})

	// Getting verses of coros
	for i, himno := range himnos {
		parrafos, _ := new(models.Parrafo).GetParrafos(himno.ID)
		himnos[i].Parrafos = parrafos

		db.Collection("Songs").InsertOne(context.TODO(), HimnoToSong(
			&himnos[i],
			songbook_result.InsertedID.(primitive.ObjectID),
			[]primitive.ObjectID{
				tema_coros.InsertedID,
			},
		))
	}
}

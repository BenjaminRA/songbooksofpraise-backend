package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BenjaminRA/himnario-backend/db/mongodb"
	"github.com/BenjaminRA/himnario-backend/db/sqlite"
	"github.com/BenjaminRA/himnario-backend/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tcolgate/mp3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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

func HimnoToSong(himno *models.Himno, songbook_id primitive.ObjectID, categories []primitive.ObjectID) models.Song {
	chords := true
	new_verses := []models.Verse{}
	for idx, parrafo := range himno.Parrafos {
		if chords {
			chords = parrafo.Acordes != ""
		}
		new_verses = append(new_verses, models.Verse{
			Text:        parrafo.Parrafo,
			Chorus:      parrafo.Coro,
			OrderNumber: idx + 1,
		})
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

	if himno.ID > 517 {
		himno.ID = himno.ID - 517
	}

	return models.Song{
		ID:           primitive.NewObjectID(),
		Number:       himno.ID,
		Title:        himno.Titulo,
		Chords:       chords,
		Verses:       new_verses,
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

func contains(array []primitive.ObjectID, element primitive.ObjectID) bool {
	for _, el := range array {
		if el == element {
			return true
		}
	}

	return false
}

func checkIfExists(db *mongo.Database, song *models.Himno) bool {
	itemCount, err := db.Collection("Songs").CountDocuments(context.TODO(), bson.M{"number": song.ID})
	if err != nil {
		panic(err)
	}

	return itemCount > 0
}

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

func addCategoryToHimno(himnos_tema *map[int]([]primitive.ObjectID), himno_id int, category_id primitive.ObjectID) {
	if _, found := (*himnos_tema)[himno_id]; !found {
		(*himnos_tema)[himno_id] = []primitive.ObjectID{}
	}
	if !contains((*himnos_tema)[himno_id], category_id) {
		(*himnos_tema)[himno_id] = append((*himnos_tema)[himno_id], category_id)
	}
}

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

	// Closing resources
	defer sqlite.Disconnect()
	defer mongodb.Disconnect()

	db := mongodb.GetMongoDBConnection()

	songbook := models.Songbook{
		ID:          primitive.NewObjectID(),
		Title:       "Himnos y Cánticos del Evangelio",
		Language:    "es",
		Description: "...",
		Country: models.Country{
			Country: "Chile",
			Code:    "CL",
		},
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

	todos_id := primitive.NewObjectID()

	temas = append(temas, models.Tema{
		InsertedID: todos_id,
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
				addCategoryToHimno(&himnos_category, temas[i].Himnos[j].ID, temas[i].InsertedID)
				temas[i].SubTemas[j].Himnos[k].Parrafos, _ = new(models.Parrafo).GetParrafos(temas[i].SubTemas[j].Himnos[k].ID)
			}
		}

	}

	for _, tema := range temas {

		db.Collection("Categories").InsertOne(context.TODO(), models.Category{
			ID:         tema.InsertedID,
			SongbookID: songbook_result.InsertedID.(primitive.ObjectID),
			Category:   tema.Tema,
		})

		children := []primitive.ObjectID{}
		if len(tema.SubTemas) > 0 {
			for _, subtema := range tema.SubTemas {
				children = append(children, subtema.InsertedID)
				db.Collection("Categories").InsertOne(context.TODO(), models.Category{
					ID:         subtema.InsertedID,
					Category:   subtema.Tema,
					SongbookID: songbook_result.InsertedID.(primitive.ObjectID),
					ParentID:   tema.InsertedID,
				})

				for _, subtema_himno := range subtema.Himnos {
					if !checkIfExists(db, &subtema_himno) {
						db.Collection("Songs").InsertOne(context.TODO(), HimnoToSong(
							&subtema_himno,
							songbook_result.InsertedID.(primitive.ObjectID),
							himnos_category[subtema_himno.ID],
						))
					}
				}
			}
		}

		for _, himno := range tema.Himnos {
			if !checkIfExists(db, &himno) {
				db.Collection("Songs").InsertOne(context.TODO(), HimnoToSong(
					&himno,
					songbook_result.InsertedID.(primitive.ObjectID),
					himnos_category[himno.ID],
				))
			}
		}
		db.Collection("Categories").UpdateByID(context.TODO(), tema.InsertedID, bson.M{
			"$set": bson.M{
				"children_ID": children,
			},
		})
	}

	for _, himno := range himnos {
		addCategoryToHimno(&himnos_category, himno.ID, todos_id)
		if !checkIfExists(db, &himno) {
			db.Collection("Songs").InsertOne(context.TODO(), HimnoToSong(
				&himno,
				songbook_result.InsertedID.(primitive.ObjectID),
				himnos_category[himno.ID],
			))
		}
	}

	db.Collection("Songbooks").InsertOne(context.TODO(), songbook)

	songbook = models.Songbook{
		ID:          primitive.NewObjectID(),
		Title:       "Coros",
		Language:    "es",
		Description: "...",
		Country: models.Country{
			Country: "Chile",
			Code:    "CL",
		},
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
		SongbookID: songbook_result.InsertedID.(primitive.ObjectID),
		Category:   tema_coros.Tema,
	})

	// Getting verses of coros
	for i, himno := range himnos {
		parrafos, _ := new(models.Parrafo).GetParrafos(himno.ID)
		himnos[i].Parrafos = parrafos
		fmt.Println(parrafos)

		db.Collection("Songs").InsertOne(context.TODO(), HimnoToSong(
			&himnos[i],
			songbook_result.InsertedID.(primitive.ObjectID),
			[]primitive.ObjectID{
				tema_coros.InsertedID,
			},
		))
	}

	// result, _ := db.Collection("Songbooks").InsertOne(context.TODO(), songbook)

	// for i := 0; i < len(temas); i++ {
	// 	temas[i].InsertedID = primitive.NewObjectID()
	// 	for j := 0; j < len(temas[i].Himnos); j++ {
	// 		// temas[i].Himnos[j] = HimnoToSong(temas[i].Himnos[j], result.InsertedID, temas[i].InsertedID)
	// 	}

	// 	for j := 0; j < len(temas[i].SubTemas); j++ {
	// 		temas[i].SubTemas[j].Himnos, _ = temas[i].SubTemas[j].GetSubTemaHimnos()
	// 		for k := 0; k < len(temas[i].SubTemas[j].Himnos); k++ {
	// 		}
	// 	}

	// }

	// cursor, err := db.Collection("Categories").Aggregate(context.TODO(), []bson.M{
	// 	{"$match": bson.M{"parent_id": primitive.Null{}}},
	// 	{"$graphLookup": bson.M{
	// 		"from":             "Categories",
	// 		"startWith":        "$parent_id",
	// 		"connectFromField": "parent_id",
	// 		"connectToField":   "_id",
	// 		"as":               "parent",
	// 	}},
	// 	{"$graphLookup": bson.M{
	// 		"from":             "Categories",
	// 		"startWith":        "$_id",
	// 		"connectFromField": "children",
	// 		"connectToField":   "_id",
	// 		"restrictSearchWithMatch": bson.M{
	// 			"children": bson.M{
	// 				"$not": bson.M{
	// 					"$size": 0,
	// 				},
	// 			},
	// 		},
	// 		"as": "sub_categories",
	// 	}},
	// 	{"$project": bson.M{"_id": 0, "category": 1, "parent": 1, "sub_categories": 1}},
	// 	{"$project": bson.M{
	// 		"categories.category": bson.M{
	// 			"$filter": bson.M{
	// 				"input": "$categories.category",
	// 				"as":    "category",
	// 				"cond": bson.M{
	// 					"$eq": bson.A{"$$category", "Todos"},
	// 				},
	// 			},
	// 		},
	// 	}},
	// })

	// if err != nil {
	// 	panic(err)
	// }

	// for cursor.Next(context.TODO()) {
	// 	var element interface{}
	// 	err := cursor.Decode(&element)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	// fmt.Println(element)
	// }

	// coros, _ := new(models.Himno).GetCoros()
	// for idx, coro := range coros {
	// 	// fmt.Println("ID:", coro.ID, "Título:", coro.Titulo)
	// 	parrafos, _ := new(models.Parrafo).GetParrafos(coro.ID)
	// 	coros[idx].Parrafos = parrafos

	// 	temas, _ := new(models.Tema).GetTemas(coro.ID)
	// 	coros[idx].Temas = temas
	// }

}

package migration

import (
	"context"
	"fmt"

	"github.com/BenjaminRA/himnario-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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
